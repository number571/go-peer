import multiprocessing, sys, requests, time, json, hashlib, yaml

HLE_URL = "" # it is being overwritten
HLT_URL = "" # it is being overwritten
FRIENDS = {} # it is being overwritten

def main():
    init_load_config("config.yml")

    parallel_run(input_task, output_task)

def init_load_config(cfg_path):
    global HLE_URL, HLT_URL, FRIENDS

    with open(cfg_path, "r") as stream:
        try:
            config_loaded = yaml.safe_load(stream)
        except yaml.YAMLError as e:
            print("@ failed load config")
            exit(2)
    
    HLT_URL = "http://" + config_loaded["hlt_host"]
    HLE_URL = "http://" + config_loaded["hle_host"]
    FRIENDS = config_loaded["friends"]

def parallel_run(*fns):
    proc = []
    for fn in fns:
        p = multiprocessing.Process(target=fn)
        p.start()
        proc.append(p)
    for p in proc:
        p.join()

def output_task():
    # GET SETTING = MESSAGES_CAPACITY
    resp_hlt = requests.get(
        HLT_URL+"/api/config/settings"
    )
    if resp_hlt.status_code != 200:
        print("@ got response error from HLT (/api/config/settings)")
        exit(1)
    
    try:
        messages_capacity = json.loads(resp_hlt.content)["messages_capacity"]
    except ValueError:
        print("@ got response invalid data from HLT (/api/config/settings)")
        exit(1)
    
    global_pointer = -1
    while True:
        # GET INITIAL POINTER OF MESSAGES
        resp_hlt = requests.get(
            HLT_URL+"/api/storage/pointer"
        )
        if resp_hlt.status_code != 200:
            print("@ got response error from HLT (/api/storage/pointer)")
            time.sleep(1)
            continue 

        try:
            pointer = int(resp_hlt.content)
        except ValueError:
            print("@ got response invalid data from HLT (/api/storage/pointer)")
            time.sleep(1)
            continue

        if global_pointer == -1:
            global_pointer = pointer

        if global_pointer == pointer:
            time.sleep(1)
            continue
    
        # GET ALL MESSAGES FROM CURRENT POINTER TO GOT POINTER
        while global_pointer != pointer:
            global_pointer = (global_pointer + 1) % messages_capacity

            resp_hlt = requests.get(
                HLT_URL+"/api/storage/hashes?id="+f"{(global_pointer - 1) % messages_capacity}"
            )
            if resp_hlt.status_code != 200:
                break 
            
            resp_hlt = requests.get(
                HLT_URL+"/api/network/message?hash="+resp_hlt.content.decode("utf8")
            )
            if resp_hlt.status_code != 200:
                break 

            # TRY DECRYPT GOT MESSAGE
            resp_hle = requests.post(
                HLE_URL+"/api/message/decrypt", 
                data=resp_hlt.content
            )
            if resp_hle.status_code != 200:
                continue 

            try:
                json_resp = json.loads(resp_hle.content)
            except ValueError:
                print("@ got response invalid data from HLE (/api/message/decrypt)")
                continue
        
            # CHECK GOT PUBLIC KEY IN FRIENDS LIST
            friend_name = ""
            user_id = hashlib.sha256(json_resp["public_key"].encode('utf-8')).hexdigest()

            friend_exist = False 
            for k, v in FRIENDS.items():
                friend_id = hashlib.sha256(v.encode('utf-8')).hexdigest()
                if user_id == friend_id:
                    friend_name = k
                    friend_exist = True
                    break

            got_data = bytes.fromhex(json_resp["hex_data"]).decode('utf-8')
            print(f"[{friend_name}]: {got_data}\n> ", end="")

def input_task():
    friend = ""

    sys.stdin = open(0)
    while True:
        msg = input("> ")
        if len(msg) == 0:
            print("@ got null message")
            continue

        if msg.startswith("/friend "):
            try:
                _friend = FRIENDS[msg[len("/friend "):].strip()]
            except KeyError:
                print("@ got invalid friend name")
                continue
            friend = _friend
            continue

        if msg.startswith("/exit"):
            print("$ goodbye")
            sys.exit(0)
        
        if friend == "":
            print("@ friend is null, use /friend to set")
            continue 

        resp_hle = requests.post(
            HLE_URL+"/api/message/encrypt", 
            json={"public_key": friend, "hex_data": msg.encode("utf-8").hex()}
        )
        if resp_hle.status_code != 200:
            print("@ got response error from HLE (/api/message/encrypt)")
            continue 
        
        resp_hlt = requests.post(
            HLT_URL+"/api/network/message", 
            data=resp_hle.content
        )
        if resp_hlt.status_code != 200:
            print("@ got response error from HLT (/api/network/message)")
            continue 

if __name__ == "__main__":
    main()
