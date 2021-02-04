import java.io.*;
import java.util.*;

import io.prometheus.client.Collector;
import io.prometheus.client.exporter.HTTPServer;

class HDFSDiskCollector extends Collector {

	String path = "/custom/prometheus_export/get_system_metrics.sh";
	String dir = "/custom/hdfs_dir";

	HDFSDiskCollector(String path, String dir) {
		this.path = path;
		this.dir = dir;
	}

	public List<MetricFamilySamples> collect() {

		String[] metricNames = { "directory_size", "disk_utilization", "disk_iops", "disk_read_mbps", "disk_write_mbps",
				"disk_read_mb", "disk_write_mb" };
		List<String> lableList = Arrays.asList("directory");
		List<String> lableValueList = Arrays.asList(dir);
		String[] helpInformations = { "Get current directory size use du",
				"Get the disk utilization which mount in this directory use df",
				"Get the disk iops which mount in this directory use iostat",
				"Get the disk read throughput which mount in this directory use iostat(mb/s)",
				"Get the disk write throughput which mount in this directory use iostat(mb/s)",
				"Get the disk read data size which mount in this directory use iostat(mb)",
				"Get the disk write data size which mount in this directory use iostat(mb)" };

		MetricFamilySamples.Sample sample;
		MetricFamilySamples samples;
		List<MetricFamilySamples> mfs = new ArrayList<MetricFamilySamples>();

		try {
			// Use the Linux System's command to get system metrics
			Process ps = Runtime.getRuntime().exec(new String[] { path, dir });
			BufferedReader br = new BufferedReader(new InputStreamReader(ps.getInputStream()));
			String metric;
			for (int index = 0; index < 7; index++) {
				metric = br.readLine();
				sample = new MetricFamilySamples.Sample(metricNames[index], lableList, lableValueList,
						Double.parseDouble(metric));
				samples = new MetricFamilySamples(metricNames[index], Type.GAUGE, helpInformations[index],
						Arrays.asList(sample));
				mfs.add(samples);
			}
		} catch (Exception e) {
			e.printStackTrace();
		}

		return mfs;
	}
}

public class HDFSDiskExport {
	public static void main(String[] args) {
		int port = args.length > 0 ? Integer.parseInt(args[0]) : 10318;
		String path = args.length > 1 ? args[1] : "/custom/prometheus_export/get_system_metrics.sh";
		String dir = args.length > 2 ? args[2] : "/custom/hdfs_dir";

		HTTPServer server = null;
		try {
			server = new HTTPServer(port);
			// register metrics collector
			new HDFSDiskCollector(path, dir).register();

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