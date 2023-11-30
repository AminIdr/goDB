output = open("commands2.txt", "w")

def set(key, val):
    return "curl -X POST -H \"Content-Type: application/json\" -d \"{\\\"key\\\": \\\"" + key + "\\\", \\\"value\\\": \\\"" + val + "\\\"}\" http://localhost:8080/set"

def delete(key):
    return "curl http://localhost:8080/del?key=" + key

def get(key):
    return "curl http://localhost:8080/get?key=" + key

for i in range(1000):
    key = str(i)
    val = str(i)
    output.write(set(key, val) + '\n')