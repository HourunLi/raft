import os

for i in range(3):
    os.system("start ./test.exe > out{}.txt".format(i))
