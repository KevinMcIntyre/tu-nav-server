# TU Nav Server

## Instructions

 - Make sure you have [Go](https://golang.org/) installed and [setup](https://golang.org/doc/install), minimum version 1.5.1
 - Setup [Godep](https://github.com/tools/godep)

####Clone the repository, enter the directory

    $ git clone https://github.com/KevinMcIntyre/tu-nav-server.git && cd tu-nav-server
    
####Pull the dependencies

    $ godep restore
    
####Compile the source

    $ go build main.go
    
####Run the executable

    $ ./main
