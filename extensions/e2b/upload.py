from e2b import Sandbox

sandbox = Sandbox(template="base")

with open("./README.md", "rb") as f:
    remote_path = sandbox.upload_file(f)  # $HighlightLine

sandbox.close()
