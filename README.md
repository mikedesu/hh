# hosthunter

Hunt for SSRFs in the host header.

# building

```
go build hh.go
```

# usage

```
Usage of ./hh:
  -f string
    	input file name
  -h string
    	collaborator or interactsh host (default "yoururl.com")
```

# examples

```
./hh -f domains_up -h my.collaborator.url 
```

