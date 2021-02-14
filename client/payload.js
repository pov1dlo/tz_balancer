var frequency
//счетчик отправленных пакетов
var counter = 0

var xhr = new XMLHttpRequest();
xhr.open("GET", "config.json")
xhr.send();

xhr.onload = function(){
    config = JSON.parse(this.responseText);
    frequency = config.frequency;
    console.log("frequency:", frequency);
    
    setInterval(postPayload, frequency);
};

//var weightMillisec = getRandonInt(10)

function getRandonInt(max){

    rand = Math.floor(Math.random() * Math.floor(max));

    if (rand == 0) {
        rand = getRandonInt(max);
    }

    return rand;

}

function Payload(){

    this.price      = getRandonInt(1000);
    this.quantity   = getRandonInt(10);
    this.amount     = this.price * this.quantity;
    this.object     = getRandonInt(10);
    this.method     = getRandonInt(10);

}

function postPayload(){

    var maxPackage  = getRandonInt(10000)
    var package     = []
    for ( i=0; i < maxPackage; i++ ){
        package.push(new Payload())
    }

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:8080/upload", true);
    xhr.setRequestHeader("Content-type", "application/json")
    xhr.send(JSON.stringify(package))

    xhr.onload = function() {
        if (xhr.status != 200) { 
            //console.log(xhr.status, xhr.statusText);
        } else { 
            counter++;
            console.log("Количество обработанных пакетов: ", counter);
            console.log(xhr.responseText);
        }
      };

}

//postPayload()
