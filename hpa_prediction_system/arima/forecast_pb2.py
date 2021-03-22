# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: forecast.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='forecast.proto',
  package='time_series_forecast',
  syntax='proto3',
  serialized_options=b'ZEgithub.com/k8s-autoscaling/hpa_prediction_system/time_series_forecast',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x0e\x66orecast.proto\x12\x14time_series_forecast\"0\n\x0f\x46orecastRequest\x12\x0c\n\x04\x64\x61ta\x18\x01 \x01(\t\x12\x0f\n\x07minutes\x18\x02 \x01(\x05\"!\n\x10\x46orecastResponse\x12\r\n\x05value\x18\x01 \x01(\x02\x32v\n\x0f\x46orecastService\x12\x63\n\x10GetForeCastValue\x12%.time_series_forecast.ForecastRequest\x1a&.time_series_forecast.ForecastResponse\"\x00\x42GZEgithub.com/k8s-autoscaling/hpa_prediction_system/time_series_forecastb\x06proto3'
)




_FORECASTREQUEST = _descriptor.Descriptor(
  name='ForecastRequest',
  full_name='time_series_forecast.ForecastRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='data', full_name='time_series_forecast.ForecastRequest.data', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='minutes', full_name='time_series_forecast.ForecastRequest.minutes', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=40,
  serialized_end=88,
)


_FORECASTRESPONSE = _descriptor.Descriptor(
  name='ForecastResponse',
  full_name='time_series_forecast.ForecastResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='time_series_forecast.ForecastResponse.value', index=0,
      number=1, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=90,
  serialized_end=123,
)

DESCRIPTOR.message_types_by_name['ForecastRequest'] = _FORECASTREQUEST
DESCRIPTOR.message_types_by_name['ForecastResponse'] = _FORECASTRESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ForecastRequest = _reflection.GeneratedProtocolMessageType('ForecastRequest', (_message.Message,), {
  'DESCRIPTOR' : _FORECASTREQUEST,
  '__module__' : 'forecast_pb2'
  # @@protoc_insertion_point(class_scope:time_series_forecast.ForecastRequest)
  })
_sym_db.RegisterMessage(ForecastRequest)

ForecastResponse = _reflection.GeneratedProtocolMessageType('ForecastResponse', (_message.Message,), {
  'DESCRIPTOR' : _FORECASTRESPONSE,
  '__module__' : 'forecast_pb2'
  # @@protoc_insertion_point(class_scope:time_series_forecast.ForecastResponse)
  })
_sym_db.RegisterMessage(ForecastResponse)


DESCRIPTOR._options = None

_FORECASTSERVICE = _descriptor.ServiceDescriptor(
  name='ForecastService',
  full_name='time_series_forecast.ForecastService',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=125,
  serialized_end=243,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetForeCastValue',
    full_name='time_series_forecast.ForecastService.GetForeCastValue',
    index=0,
    containing_service=None,
    input_type=_FORECASTREQUEST,
    output_type=_FORECASTRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_FORECASTSERVICE)

DESCRIPTOR.services_by_name['ForecastService'] = _FORECASTSERVICE

# @@protoc_insertion_point(module_scope)