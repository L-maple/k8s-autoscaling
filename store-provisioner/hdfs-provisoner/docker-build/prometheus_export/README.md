# Dependence

This code compile and run on Java 8, but I think it can be compile and run on any version of Java.

The package of Java you can download in [Oracle](https://www.oracle.com/java/technologies/javase/javase8-archive-downloads.html).(Because you must accept the license manually, you can search `jdk-8u11-linux-x64.tar.gz` for Linux and `jdk-8u11-windows-x64.exe` for Windows)

# Configure on Linux

Assuming that we download the installation package in `/home/ubuntu`, you can configure Java by this.

```shell
tar -xf jdk-8u11-linux-x64.tar.gz
cat << EOF >> /etc/profile
export JAVA_HOME=/home/ubuntu/jdk1.8.0_11
export JRE_HOME=\$JAVA_HOME/jre
export CLASSPATH=.:\$CLASSPATH:\$JAVA_HOME/lib:\$JAVA_HOME/jre/lib
export PATH=\$PATH:\$JAVA_HOME/bin:\$JAVA_HOME/jre/bin
EOF
source /etc/profile
```

Then, you can compile the code by `javac -cp .:lib/* DirectorySizeExport.java` and run `java -cp .:lib/* DirectorySizeExport [port] [path] [dir]`(the usage can refer [here](../../))

# Configure on Windows

The install of Java and Eclipse you can finish easy, just click the execute file. Configure the environment variable can refer [here](https://jingyan.baidu.com/article/8275fc86b2cf7b46a03cf6bc.html), the install package of Eclipse can download [here](http://mirrors.neusoft.edu.cn/eclipse/technology/epp/downloads/release/oxygen/3a/eclipse-java-oxygen-3a-win32-x86_64.zip).([refer](https://www.jianshu.com/p/b8d7c5438302))

You should import all jar file of this lib directory.

1. Right click your project
2. Click Build Path -> Add External Archives
3. Select all jar file you download

Now, you can write Prometheus Java Client Code.


