import os
import sys
import threading

# https://codereview.stackexchange.com/a/17959
class LogPipe(threading.Thread):

    def __init__(self):
        """Setup the object with a logger and a loglevel
        and start the thread
        """
        threading.Thread.__init__(self)
        self.daemon = False
        self.fdRead, self.fdWrite = os.pipe()
        self.pipeReader = os.fdopen(self.fdRead, encoding='utf-8', errors='ignore')
        self.start()

    def fileno(self):
        """Return the write file descriptor of the pipe
        """
        return self.fdWrite

    def run(self):
        """Run the thread, logging everything.
        """
        while line := self.pipeReader.readline():
            sys.stdout.write(line)

        self.pipeReader.close()

    def close(self):
        """Close the write end of the pipe.
        """
        os.close(self.fdWrite)
