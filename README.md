<h1 align="center">
    <br>
    <a href="https://github.com/aureleoules/gocaml"><img src="https://raw.githubusercontent.com/aureleoules/gocaml/master/assets/icon.png" alt="GOCAML" width="200"></a>
    <br>
    GOCAML
    <br>
    <a href="https://travis-ci.org/aureleoules/gocaml"><img src="https://travis-ci.org/aureleoules/gocaml.svg?branch=master"></a>  
    <br>
</h1>

<h4 align="center">Evaluate CAML on Discord</h4>

<img src="https://raw.githubusercontent.com/aureleoules/gocaml/master/assets/example.png" alt="example">

## Install
* `go get github.com/aureleoules/gocaml`
* Edit `.env` and add your Discord BOT Token
* `go build`

## Get started
To evaluate CAML code on your discord, simply tag your Discord BOT, and enter your code in a triple backtick block as follows:

> @GOCAML#1234  
> \`\`\`ocaml  
> let a = 2;;  
> \`\`\`  
  
will result in  

```ocaml
OCaml version 4.02.3

# val a : int = 2
```  

## License

[MIT](https://github.com/aureleoules/glaze/blob/master/LICENSE) © [Aurèle Oulès](https://www.aureleoules.com)