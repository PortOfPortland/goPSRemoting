package goPSRemoting 

import (
        "os/exec"
        "bytes"
        "runtime"
        "strings"
        "errors"
)

func runCommand(args ...string) (string, error) {
        cmd := exec.Command(args[0], args[1:]...)

        var out bytes.Buffer
        var err bytes.Buffer

        cmd.Stdout = &out 
        cmd.Stderr = &err
        cmd.Run()

        // convert err to an error type if there is an error returned
        var e error
        if err.String() != "" {
                e = errors.New(err.String())
        }

        return strings.TrimRight(out.String(), "\r\n"), e
}

func RunPowershellCommand(username string, password string, server string, command string, usessl string, usessh string) (string, error) {
        var pscommand string

        if runtime.GOOS == "windows" {
                pscommand = "powershell.exe"
        } else {
                pscommand = "pwsh"
        }

        var winRMPre string

        if (usessh == "1") {
                winRMPre = "$s = New-PSSession -HostName " + server + " -Username " + username + " -SSHTransport"
        } else {
                winRMPre = "$SecurePassword = '" + password + "' | ConvertTo-SecureString -AsPlainText -Force; $cred = New-Object System.Management.Automation.PSCredential -ArgumentList '" + username + "', $SecurePassword; $s = New-PSSession -ComputerName " + server + " -Credential $cred"
        }

        var winRMPost string

        if runtime.GOOS == "windows" {
                winRMPost = "; Invoke-Command -Session $s -Scriptblock { " + command + " }; Remove-PSSession $s"
        } else {
                winRMPost = "; Invoke-Command -Session $s -Scriptblock { powershell '" + command + "' }; Remove-PSSession $s"
        }

        var winRMCommand string

        if (usessl == "1") {
                winRMCommand = winRMPre + " -UseSSL" + winRMPost
        } else {
                winRMCommand = winRMPre + winRMPost
        }

        out, err := runCommand(pscommand, "-command", winRMCommand) 

        return out, err
}