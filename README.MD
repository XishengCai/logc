### golang  logrus 学习

#### logrus特性：

   1. 完全兼容golang标准库日志模块：logrus拥有六种日志级别：debug、info、warn、error、fatal和panic，这是golang标准库日志模块的API的超集。如果您的项目使用标准库日志模块，完全可以以最低的代价迁移到logrus上。
   2. 可扩展的Hook机制：允许使用者通过hook的方式将日志分发到任意地方，如本地文件系统、标准输出、logstash、elasticsearch或者mq等，或者通过hook定义日志内容和格式等。
   3. 可选的日志输出格式：logrus内置了两种日志格式，JSONFormatter和TextFormatter，如果这两个格式不满足需求，可以自己动手实现接口Formatter，来定义自己的日志格式。
   4. Field机制：logrus鼓励通过Field机制进行精细化的、结构化的日志记录，而不是通过冗长的消息来记录日志。
   5. logrus是一个可插拔的、结构化的日志框架

### hook
   - logrus的hook接口定义如下，其原理是每此写入日志时拦截，修改logrus.Entry。
   ```
   // logrus在记录Levels()返回的日志级别的消息时会触发HOOK，
   // 按照Fire方法定义的内容修改logrus.Entry。
   type Hook interface {
       Levels() []Level
       Fire(*Entry) error
   }
   ```

   - 一个简单自定义hook如下，DefaultFieldHook定义会在所有级别的日志消息中加入默认字段appName="myAppName"。
   ```
   type DefaultFieldHook struct {
   }

   func (hook *DefaultFieldHook) Fire(entry *log.Entry) error {
       entry.Data["appName"] = "MyAppName"
       return nil
   }

   func (hook *DefaultFieldHook) Levels() []log.Level {
       return log.AllLevels
   }
   ```

   在初始化前调用log.AddHook(hook)添加相应的hook即可。

### 缺点

   尽管logrus有诸多优点，但是为了灵活性和可扩展性，官方也削减了很多实用的功能，例如：
   没有提供行号和文件名的支持
   输出到本地文件系统没有提供日志分割功能
   官方没有提供输出到ELK等日志处理中心的功能

### 记录文件名和行号
   logrus的一个很致命的问题就是没有提供文件名和行号，这在大型项目中通过日志定位问题时有诸多不便。Github上的logrus的
   issue#63：Log filename and line number创建于2014年，四年过去了仍是open状态~~~
   网上给出的解决方案分位两类，一就是自己实现一个hook；二就是通过装饰器包装logrus.Entry。
   两种方案网上都有很多代码，但是大多无法正常工作。但总体来说，解决问题的思路都是对的：
   通过标准库的runtime模块获取运行时信息，并从中提取文件名，行号和调用函数名。

   标准库runtime模块的Caller(skip int)函数可以返回当前goroutine调用栈中的文件名，
   行号，函数信息等，参数skip表示表示返回的栈帧的层次，0表示runtime.Caller的调用着。
   返回值包括响应栈帧层次的pc(程序计数器)，文件名和行号信息。为了提高效率，
   我们先通过跟踪调用栈发现，从runtime.Caller()的调用者开始，到记录日志的生成代码之间
   ，大概有8到11层左右，所有我们在hook中循环第8到11层调用栈应该可以找到日志记录的生产代码。

   此外，runtime.FuncForPC(pc uintptr) *Func可以返回指定pc的函数信息。
   所有我们要实现的hook也是基于以上原理，使用runtime.Caller()依次循环调用栈的第7~11层，
   过滤掉sirupsen包内容，那么第一个非siupsenr包就认为是我们的生产代码了，
   并返回pc以便通过runtime.FuncForPC()获取函数名称。然后将文件名、行号和函数名组装为source字段塞到logrus.Entry中即可。

### 日志本地文件分割
   logrus本身不带日志本地文件分割功能，但是我们可以通过file-rotatelogs进行日志本地文件分割。
   每次当我们写入日志的时候，logrus都会调用file-rotatelogs来判断日志是否要进行切分。
   ```
   import (
       "github.com/lestrrat-go/file-rotatelogs"
       "github.com/rifflock/lfshook"
       log "github.com/sirupsen/logrus"
       "time"
   )

   func newLfsHook(logLevel *string, maxRemainCnt uint) log.Hook {
       writer, err := rotatelogs.New(
           logName+".%Y%m%d%H",
           // WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
           rotatelogs.WithLinkName(logName),

           // WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
           rotatelogs.WithRotationTime(time.Hour),

           // WithMaxAge和WithRotationCount二者只能设置一个，
           // WithMaxAge设置文件清理前的最长保存时间，
           // WithRotationCount设置文件清理前最多保存的个数。
           //rotatelogs.WithMaxAge(time.Hour*24),
           rotatelogs.WithRotationCount(maxRemainCnt),
       )

       if err != nil {
           log.Errorf("config local file system for logger error: %v", err)
       }

       level, ok := logLevels[*logLevel]

       if ok {
           log.SetLevel(level)
       } else {
           log.SetLevel(log.WarnLevel)
       }

       lfsHook := lfshook.NewHook(lfshook.WriterMap{
           log.DebugLevel: writer,
           log.InfoLevel:  writer,
           log.WarnLevel:  writer,
           log.ErrorLevel: writer,
           log.FatalLevel: writer,
           log.PanicLevel: writer,
       }, &log.TextFormatter{DisableColors: true})

       return lfsHook
   }
   ```

   [抄袭博客](https://blog.csdn.net/wslyk606/article/details/81670713)