import os
import time

from e2b.sandbox import Sandbox
from e2b.sandbox.process import ProcessMessage

watcher = None

def create_watcher(sandbox):  
    # Start filesystem watcher for the /home directory
    watcher = sandbox.filesystem.watch_dir("/home")  
    watcher.add_event_listener(lambda event: print(event))  
    watcher.start()  

def on_stdout(message: ProcessMessage):
    print(f"Process {message.process_id} stdout: {message.data}")


sandbox = Sandbox(template="base")

create_watcher(sandbox)  

sandbox.keep_alive(101)
# Create files in the /home directory inside the playground
# We'll receive notifications for these events through the watcher we created above.
start = time.time()
for i in range(10):
    # `filesystem.write()` will trigger two events:
    # 1. 'Create' when the file is created
    # 2. 'Write' when the file is written to
    sandbox.filesystem.write(f"/home/file{i}.txt", f"Hello World {i}!")
    sandbox.filesystem.make_dir(f"/home/dir{i}")
    _ = sandbox.filesystem.read(f"/home/file{i}.txt")
    _ = sandbox.filesystem.list(f"/home/")
    _ = sandbox.filesystem.read_bytes(f"/home/file{i}.txt")
    sandbox.filesystem.write_bytes(f"/home/file{i}.txt", b"Hello World {i}!")
    sandbox.filesystem.remove(f"/home/file{i}.txt")
    sandbox.filesystem.remove(f"/home/dir{i}")

end = time.time()
print(f"Time taken: {end - start}")

# now doing the same thing with the fs api
start = time.time()
for i in range(10):
    os.mkdir(f"dir{i}")
    with open(f"file{i}.txt", "w") as f:
        _ = f.write(f"Hello World {i}!")
    with open(f"file{i}.txt", "rb") as f:
        data = f.read()
    os.remove(f"file{i}.txt")
    os.rmdir(f"dir{i}")
end = time.time()
print(f"Time taken: {end - start}")

# now doing the same thing with the process api
start = time.time()
for i in range(10): 
    proc = sandbox.process.start(cmd=f"cat 'Hello World {i}!' > file{i}.txt", on_stdout=on_stdout)
    proc.wait()
    proc = sandbox.process.start(cmd=f"cat file{i}.txt")
    proc.wait()
    proc = sandbox.process.start(cmd=f"ls")
    proc.wait()
    proc = sandbox.process.start(cmd=f"cat file{i}.txt", on_stdout=on_stdout)
    proc.wait()
    sandbox.process.start(cmd=f"rm file{i}.txt")
end = time.time()
print(f"Time taken: {end - start}")
    

sandbox.close()
