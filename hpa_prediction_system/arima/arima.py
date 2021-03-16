from concurrent import futures
import os
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from datetime import timedelta
from statsmodels.tsa.arima_model import ARIMA
from sklearn.metrics import mean_squared_error
from scipy.ndimage.filters import gaussian_filter
import matplotlib.dates as mdates
import warnings
import random
import forecast_pb2
import grpc
import forecast_pb2_grpc


warnings.filterwarnings("ignore")
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '2'


def arima_model(series, data_split, params, future_periods, log):
    # log transformation of data if user selects log as true
    if log:
        series_dates = series.index
        series = pd.Series(np.log(series), index=series.index)

    # create training and testing data sets based on user split fraction
    size = int(len(series) * data_split)
    train, test = series[0:size], series[size:len(series)]
    history = [val for val in train]
    predictions = []

    # creates a rolling forecast by testing one value from the test set, and then add that test value
    # to the model training, followed by testing the next value in the series
    for t in range(len(test)):
        model = ARIMA(history, order=(params[0], params[1], params[2]))
        model_fit = model.fit(disp=0)
        output = model_fit.forecast()
        yhat = output[0]
        predictions.append(yhat[0])
        obs = test[t]
        history.append(obs)

    # forecasts future periods past the input testing series based on user input
    model = ARIMA(history, order=(params[0], params[1], params[2]))
    model_fit = model.fit(disp=0)
    future_forecast = model_fit.forecast(future_periods)[0]
    future_dates = [test.index[-1] + timedelta(i * 365 / 12) for i in range(1, future_periods+1)]
    test_dates = test.index

    # if the data was originally log transformed, the inverse transformation is performed
    if log:
        predictions = np.exp(predictions)
        test = pd.Series(np.exp(test), index=test_dates)
        future_forecast = np.exp(future_forecast)

    # creates pandas series with datetime index for the predictions and forecast values
    forecast = pd.Series(future_forecast, index=future_dates)
    predictions = pd.Series(predictions, index=test_dates)

    # generating plots to compare the prediction for out of sample data to the actual test values
    fig = plt.figure()
    ax = fig.add_subplot(111)
    myFmt = mdates.DateFormatter('%m/%y')
    ax.xaxis.set_major_formatter(myFmt)
    plt.plot(predictions, c='red')
    plt.plot(test)
    plt.show()

    # calculate root mean squared errors (RMSEs) for the out-of-sample predictions
    error = np.sqrt(mean_squared_error(predictions, test))
    print('Test RMSE: %.3f' % error)

    return predictions, test, forecast


def get_ts_from_csv(csv_file):
    data = pd.read_csv(csv_file, parse_dates=['Date'], index_col="Date")
    data_monthly = data.resample('M').mean()
    data_ts = data_monthly.Close

    return data_ts


class ForecastService (
    forecast_pb2_grpc.ForecastServiceServicer
):
    def GetForeCastValue(self, request, context):
        if request.data == "":
            return forecast_pb2.ForecastResponse(value=-1)
        else:
            return forecast_pb2.ForecastResponse(value=1)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    forecast_pb2_grpc.add_ForecastServiceServicer_to_server(
        ForecastService(), server
    )

    server.add_insecure_port("[::]:50000")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    print("arima server starting...")
    serve()
    print("arima server terminated...")
