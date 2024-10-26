import os 

line = "conneroisu/groq-go"
replace = "conneroisu/groq-go/pkg/schema"
# for each json file in the directory
for file in os.listdir():
    if file.endswith(".json"):
        with open(file, "r") as f:
            data = f.read()
        data = data.replace(line, replace)
        with open(file, "w") as f:
            _ = f.write(data)
        
