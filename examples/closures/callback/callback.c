//go:build ignore

#include <math.h>

typedef float (*Callback)(float f);

float Invoke(Callback c)
{
   return c(M_PI);
}
