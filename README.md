# CLI 命令行实用程序开发基础



## 实验要求
按[文档](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)使用 selpg 章节要求测试你的程序

-----



## 设计说明
由于对go语言的理解不够深入，所以只好无奈把本次实验当成大型翻译作业，由selpg.c源码暴力转成go语言。
一个函数一个函数对应，一句一句对应，可以说是很机械了。这大概就是亲自解析命令行吧。虽然阅读了15331344同学的作业，但是还是不太熟悉很多用法。啊，生活。


-----


## 使用与测试结果
> 1.-s1 -e1 input.txt

![1-1](http://img.blog.csdn.net/20171017213607321)![1-2](http://img.blog.csdn.net/20171017213612288)

> 2.-s1 -e1 < input.txt

截部分图，说明命令可以执行，下半部分全是line...-line72
![2](http://img.blog.csdn.net/20171017214235577)


> 3.dir | selpg.exe -s1 -e1



(忽略乱码)

![3](http://img.blog.csdn.net/20171017221630366)




> 4.-s1 -e1 input.txt >output.txt


![4](http://img.blog.csdn.net/20171017214229349)



> 5.-s1 -e1 input.txt 2>error.txt


![5](http://img.blog.csdn.net/20171017215529462)





> 6.-s1 -e1 input.txt >output.txt 2>error.txt


![6](http://img.blog.csdn.net/20171017221057865)



> 不输入参数

![nopara](http://img.blog.csdn.net/20171017220042825)



-----

## 遇到的问题

![problem](http://img.blog.csdn.net/20171017220510504)

以前碰到这种问题都只是warning，不会影响编译。但是go好像很奇怪。

> 解决办法：多加了类似以下的几行无用的代码
```
s1 = s1 
crc = crc  
c = c  
line = line  
inbuf = inbuf  
```




