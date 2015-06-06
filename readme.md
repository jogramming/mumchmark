#mumchmark

Simple mumble stress testing tool

###Example usage:

    ./mumchmark -clients 10 -addr jonas747.com:64738
    What do you want to do?
    [q] Quit?
    [a] Play some audio (plays audio.mp3)
    [t] Send a text message
    a
    How many clients should play audio.mp3? leave empty for 10
    Input a number: 


Will spawn 10 clients connected to the address above

###Features:
 
 - Send text messages from a selected number of clients
 - Play audio

###Requirements:

 - To play audio you will need ffmpeg installed 