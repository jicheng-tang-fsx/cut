# How to use

```
./cut -start "05/18/2024 02:18:20" -end "05/18/2024 02:18:40" -file oms_20240517.log > my1.log
```

## Compare with C++ version
https://github.com/jicheng-tang-fsx/cut2 


| Program  | Go  | C++  |
| ---- | ---- | ---- |
| cost (second)  | 10.641 | 13.109 |


In most IO-frequent scenarios,   
Go is as fast as C++, or even faster than C++.   

std::string and std::iostream in C++ are terrible designs.   
std::string can't even handle UTF-8.   

This is just a simple small program for cutting logs.    
The C++ standard library is even 30% slower than the Go version.    


```bash
# Go version
➜  cut git:(main) time ./cut -start "05/18/2024 02:18:20" -end "05/18/2024 02:18:40"  -file /home/jicheng.tang/work/v8/oms_20240517.log > my1.log
./cut -start "05/18/2024 02:18:20" -end "05/18/2024 02:18:40" -file  > my1.log  9.26s user 1.80s system 103% cpu 10.641 total


# C++ version
➜  build git:(main) time ./cut2 -start "05/18/2024 02:18:20" -end "05/18/2024 02:18:40"  -file /home/jicheng.tang/work/v8/oms_20240517.log > my2.log
./cut2 -start "05/18/2024 02:18:20" -end "05/18/2024 02:18:40" -file  > my2.log  11.97s user 1.12s system 99% cpu 13.109 total


## check result, result is same
➜  build git:(main) diff my2.log /home/jicheng.tang/work/cut/my1.log
➜  build git:(main)
```