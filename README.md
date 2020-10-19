# code_compare_process
代码对拍平台 关于 运行代码与代码对拍的部分
此版是重构后的版本，相较于之前的，把各逻辑更明显的分离了，比之前的稍微规范些。


代码对拍评测平台 用go写的守护进程部分
start.go开始 创建缓冲为2的channel控制线程（这里go中的协程对开发者来说是线程）数量。
这附带用到的两个数据库表的sql结构文件

其中下面运行部分都是用docker-client去运行的，运行一次创建一个容器，运行完就停止容器并删除容器。

去数据库中读取version为0的数据，表示未进行对拍的记录。
之后将根据最大组数循环去跑 随机生成输入数据 的c++代码，将输出结果存到文件夹中，之后将用作对拍的输入数据。

然后循环将两个c++代码去根据对应组数的输入数据做对拍。将两个c++代码做对拍的输出结果以及运行结果（时间消耗/空间消耗）做json格式存到数据库中。
如果对拍过程中发现某一组对拍的两个代码的输出结果对比不同，那么说明对拍在这一组输入数据下的运行结果是不同的，是对拍失败的，这里就存储导致对拍失败的输入数据文本，并结束后面的对拍。

其中怎么获取c++的运行结果呢？
这里具体可以看judgepro.cpp中的代码，在go中也是开启容器后去执行这个cpp生成的可执行文件的，其中传入的参数可以参考创建容器部分传入的参数。
judgepro.cpp中的代码是fork出子进程，子进程去跑对拍的c++代码生成的可执行文件，当子进程结束后，父进程利用wait4系统调用获取子进程的运行结果。
这里的judgepro.cpp代码来源：https://blog.csdn.net/qq_38089964/article/details/86030149
在跑这个go工程之前，需要先对judgepro.cpp进行编译，因为工程中需要用到的是judgepro.cpp的可执行文件：g++ -o judgepro judgepro.cpp


把结果都存储到数据库中后，便可调用接口将这些数据做可视化处理来对比代码对拍的结果。

上面总的逻辑可以将go部分当作业务代码对数据进行处理等，且go中用docker做容器化，将cpp编译后的可执行文件挂载到容器中，可以将此可执行文件当作一个启动器。

这里主要是加强对go语言的熟悉，所以把代码对拍平台中守护进程部分重构了一下，写的不规范，多数地方做得不够好望大家见谅。

想去体验一下的可以看看：xtdlb.cn 账号密码都是development        (chrome)




现在运行代码部分也可用了，运行代码部分也放入了docker中运行。
运行代码部分是用php与go结合起来编写的。其逻辑：
当点击运行时，php将输入数据存入文件中，并将此文件与用户id与codeid存入另一个代码运行表中。同时会生成一个唯一标识id与这条记录组合。

go也会有一个守护进程去读取代码运行表，将其中未运行的用和对拍差不多的docker运行c++跑用户代码的逻辑跑出来运行结果与运行信息。再将运行结果与运行信息存入那条记录中。

在js中，用户点击运行后，ajax带着之前说的唯一标识id去查询运行结果与运行信息，然后获得返回数据后插入到html中。


