#include <stdio.h>

int main(int argc, char* argv[])
{
	if (argc > 1)
	{
		for (int i = 1; i < argc; i++)
		{
			printf("Hello, %s!\n", argv[i]);
		}
	}
	else
	{
		printf("Hello, world!\n");
	}
	
	return 0;
}