This is a small example showing libffi's closure API. The required C library must be compiled in advance:
```sh
gcc -shared -o libcallback.so -fPIC callback.c
```
