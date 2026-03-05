/*
 * spawn_helper.c — JNI native helper for spawning child processes
 * that inherit file descriptors (needed for tun2socks to access TUN fd).
 *
 * Android's ProcessBuilder closes all fds except 0/1/2 before exec().
 * This helper uses raw fork()+exec() which preserves open fds.
 */
#include <jni.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <sys/wait.h>
#include <errno.h>

/*
 * Fork and exec a process, keeping all file descriptors open.
 * Returns int[2]: { childPid, pipeReadFd }
 *   - childPid:    PID of the new process (>0), or -1 on error
 *   - pipeReadFd:  fd to read child's stdout+stderr from parent
 */
JNIEXPORT jintArray JNICALL
Java_com_jhopanstore_vpn_core_Tun2socksManager_nativeStartProcess(
        JNIEnv *env, jobject thiz,
        jobjectArray cmdArray) {

    int argc = (*env)->GetArrayLength(env, cmdArray);
    if (argc == 0) return NULL;

    /* Create a pipe so parent can read child's stdout/stderr */
    int pipefd[2];
    if (pipe(pipefd) < 0) return NULL;

    /* Build argv from Java String[] */
    char **argv = (char **) calloc(argc + 1, sizeof(char *));
    if (!argv) { close(pipefd[0]); close(pipefd[1]); return NULL; }

    for (int i = 0; i < argc; i++) {
        jstring jstr = (jstring)(*env)->GetObjectArrayElement(env, cmdArray, i);
        const char *cstr = (*env)->GetStringUTFChars(env, jstr, NULL);
        argv[i] = strdup(cstr);
        (*env)->ReleaseStringUTFChars(env, jstr, cstr);
        (*env)->DeleteLocalRef(env, jstr);
    }
    argv[argc] = NULL;

    pid_t pid = fork();

    if (pid < 0) {
        /* fork() failed */
        close(pipefd[0]);
        close(pipefd[1]);
        for (int i = 0; i < argc; i++) free(argv[i]);
        free(argv);
        return NULL;
    }

    if (pid == 0) {
        /* ---- Child process ---- */
        /* All file descriptors WITHOUT FD_CLOEXEC are inherited, including TUN fd! */

        close(pipefd[0]);   /* close read-end of pipe */

        /* Redirect stdout & stderr to the pipe so parent can read logs */
        dup2(pipefd[1], STDOUT_FILENO);
        dup2(pipefd[1], STDERR_FILENO);
        close(pipefd[1]);

        /* Reset signals */
        signal(SIGPIPE, SIG_DFL);

        /* Execute the binary */
        execv(argv[0], argv);

        /* If execv returns, it failed */
        _exit(127);
    }

    /* ---- Parent process ---- */
    close(pipefd[1]);   /* close write-end of pipe */

    for (int i = 0; i < argc; i++) free(argv[i]);
    free(argv);

    /* Return [pid, pipeReadFd] */
    jintArray result = (*env)->NewIntArray(env, 2);
    jint vals[2] = { (jint) pid, (jint) pipefd[0] };
    (*env)->SetIntArrayRegion(env, result, 0, 2, vals);
    return result;
}

/*
 * Send SIGTERM and reap the child process.
 */
JNIEXPORT void JNICALL
Java_com_jhopanstore_vpn_core_Tun2socksManager_nativeStopProcess(
        JNIEnv *env, jobject thiz,
        jint pid) {
    if (pid > 0) {
        kill((pid_t) pid, SIGTERM);
        /* Reap to prevent zombie; non-blocking first, then blocking briefly */
        int status;
        if (waitpid((pid_t) pid, &status, WNOHANG) == 0) {
            usleep(200000);  /* 200ms grace period */
            if (waitpid((pid_t) pid, &status, WNOHANG) == 0) {
                kill((pid_t) pid, SIGKILL);
                waitpid((pid_t) pid, &status, 0);
            }
        }
    }
}

/*
 * Blocking wait for child to exit.  Returns exit code.
 */
JNIEXPORT jint JNICALL
Java_com_jhopanstore_vpn_core_Tun2socksManager_nativeWaitProcess(
        JNIEnv *env, jobject thiz,
        jint pid) {
    if (pid <= 0) return -1;
    int status = 0;
    waitpid((pid_t) pid, &status, 0);
    if (WIFEXITED(status)) return WEXITSTATUS(status);
    if (WIFSIGNALED(status)) return -(WTERMSIG(status));
    return status;
}

/*
 * Check if a process is still alive.
 */
JNIEXPORT jboolean JNICALL
Java_com_jhopanstore_vpn_core_Tun2socksManager_nativeIsAlive(
        JNIEnv *env, jobject thiz,
        jint pid) {
    if (pid <= 0) return JNI_FALSE;
    return (kill((pid_t) pid, 0) == 0) ? JNI_TRUE : JNI_FALSE;
}
