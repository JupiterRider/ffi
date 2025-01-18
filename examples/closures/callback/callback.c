//go:build ignore

typedef void (*Callback)();

void Invoke(Callback c)
{
   c();
}
