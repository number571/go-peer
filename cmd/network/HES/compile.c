#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#ifdef _WIN32
	#include <windows.h>
#endif

int main(void) {
	const char *compiler = "go build";
	const char *filenames[] = {"client", "server"};
	const char *platforms[] = {"windows", "linux"};
	const char *archs[] = {"amd64", "386", "arm"};
	const char *files[] = {
		"client.go cdatabase.go cmodels.go csessions.go gconsts.go",
		"server.go sconfig.go sdatabase.go gconsts.go"
	};

	char command[BUFSIZ];
	char retfile[BUFSIZ];

	for (int k = 0; k < sizeof(filenames)/sizeof(filenames[0]); ++k) {
		for (int i = 0; i < sizeof(platforms)/sizeof(platforms[0]); ++i) {
			#ifdef _WIN32
				SetEnvironmentVariable("GOOS", platforms[i]);
				printf("set GOOS=%s\n", platforms[i]);
			#endif
			for (int j = 0; j < sizeof(archs)/sizeof(archs[0]); ++j) {
				#ifdef _WIN32
					SetEnvironmentVariable("GOARCH", archs[j]);
					printf("set GOARCH=%s\n", archs[j]);
				#endif
				snprintf(retfile, BUFSIZ, "%s_%s_%s", filenames[k], platforms[i], archs[j]);
				if (strcmp(platforms[i], "windows") == 0) {
					strcat(retfile, ".exe");
				}
				#ifdef _WIN32
					snprintf(command, BUFSIZ, "%s -o %s %s", compiler, retfile, files[k]);
				#else
					snprintf(command, BUFSIZ, "GOOS=%s GOARCH=%s %s -o %s %s", platforms[i], archs[j], compiler, retfile, files[k]);
				#endif
				printf("%s\n", command);
				system(command);
			}
		}
	}

	return 0;
}
