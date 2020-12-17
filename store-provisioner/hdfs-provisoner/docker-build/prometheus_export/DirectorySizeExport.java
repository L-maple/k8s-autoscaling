import java.io.*;
import java.util.*;

import io.prometheus.client.Collector;
import io.prometheus.client.exporter.HTTPServer;

class DirectorySizeCollector extends Collector {

	String dir = "/";
	// The different OS have different path of "du"
	String path = "/bin/du";
	DirectorySizeCollector(String path, String dir){
		this.dir = dir;
		this.path = path;
	}
	
	long ShellDus() {
		String res = "0";
		try {
			// Use the Linux System's command to get directory size
			// notice: the "du" command path may not in /bin/du
			Process ps = Runtime.getRuntime().exec(new String[] { path, "-s", dir });
			BufferedReader br = new BufferedReader(new InputStreamReader(ps.getInputStream()));
			String line;
			// Skip the previous result that maybe wrong information
			while ((line = br.readLine()) != null) {
				res = line;
			}
			// Get the size value of directory(remove the directory name)
			res = res.split("\t")[0];

		} catch (Exception e) {
			e.printStackTrace();
		}
		// return the result, the unit is KB
		return Long.parseLong(res);
	}

	public List<MetricFamilySamples> collect() {

		// Get metrics
		long len = ShellDus();

		// Build metrics value
		String metricName = "directory_size";
		List<String> lableList = Arrays.asList("directory");
		List<String> lableValueList = Arrays.asList(dir);
		MetricFamilySamples.Sample sample = new MetricFamilySamples.Sample(metricName, lableList, lableValueList, len);

		// Build metrics
		String helpInformation = "Get current directory size";
		MetricFamilySamples samples = new MetricFamilySamples(metricName, Type.GAUGE, helpInformation,
				Arrays.asList(sample));

		// add the metrics
		List<MetricFamilySamples> mfs = new ArrayList<MetricFamilySamples>();
		mfs.add(samples);
		return mfs;
	}
}

public class DirectorySizeExport {
	public static void main(String[] args) {
		int port = args.length>0 ? Integer.parseInt(args[0]) : 10318;
		String path = args.length>1 ? args[1] : "/bin/du";
		String dir = args.length>2 ? args[2] : "/";
		
		HTTPServer server = null;
		try {
			server = new HTTPServer(port);
			// register metrics collector
			new DirectorySizeCollector(path, dir).register();
			
			// As demon process, wait Prometheus to pull metrics
			while (true) {
				Thread.sleep(120000);
			}
		} catch (Exception e) {
			e.printStackTrace();
		} finally {
			if (server != null) {
				server.stop();
			}
		}
	}
}
