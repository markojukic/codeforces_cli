#include <errno.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/resource.h>
#include <unistd.h>
#include <unistd.h>

/*
Executor is called with three command line arguments:
 1. memory_limit in bytes (int)
 2. time_limit in seconds (double)
 3. command to execute

Limits are imposed using setrlimit system call, setting both soft and hard
limits before executing the provided command via exec. If the limits are
exceeded, process will be sent a SIGKILL signal.
*/
int main(int argc, char *argv[]) {
    int memory_limit = atoi(argv[1]);
    double time_limit = atof(argv[2]);
    argv += 3;
    struct rlimit limit = {
        .rlim_cur = memory_limit,
        .rlim_max = memory_limit,
    };
    if (setrlimit(RLIMIT_AS, &limit) == -1) {
        perror("setrlimit");
        return 1;
    }
    limit.rlim_cur = -1;
    limit.rlim_cur = limit.rlim_max = ceil(time_limit);
    if (setrlimit(RLIMIT_CPU, &limit) == -1) {
        perror("setrlimit");
        return 1;
    }
    if (execv(argv[0], argv) == -1) {
        perror("exec");
        return 1;
    }
    return 1;
}
