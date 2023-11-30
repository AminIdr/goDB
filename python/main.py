from random import randint

output = open("commands.txt", "w")

def set(key, val):
    return "curl -X POST -H \"Content-Type: application/json\" -d \"{\\\"key\\\": \\\"" + key + "\\\", \\\"value\\\": \\\"" + val + "\\\"}\" http://localhost:8080/set"

def delete(key):
    return "curl http://localhost:8080/del?key=" + key

for i in range(20):
    key = "amine" + str(i)
    val = "idrissi"
    output.write(set(key, val) + '\n')