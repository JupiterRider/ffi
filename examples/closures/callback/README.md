This is a small example showing libffi's closure API.
It basically just creates a callback function and passes it as an argument to another function that invokes the callback.

The required C library must be compiled in advance:
```sh
gcc -shared -o libcallback.so -fPIC callback.c
```
