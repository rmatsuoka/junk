#include <stdio.h>
#include <stdlib.h>

#include "x.h"

#define BSIZE 4096


int main(int argc, char **argv)
{
	int n;
	char buf[BSIZE];
	n = get("https://rmatsuoka.org", BSIZE-1, buf);
	if (n < 0) {
		return 1;
	}
	buf[n] = '\0';
	printf("%s\n", buf);
	return 0;
}