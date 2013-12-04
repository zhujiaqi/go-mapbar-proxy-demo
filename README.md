### go-mapbar-proxy-demo

Sends cells info, gets latlon back.

A Go-lang practice.

### Caution

Before run this, you have to register the server's IP on [http://open.mapbar.com/API_internet.jsp](http://open.mapbar.com/API_internet.jsp "Mapbar API"), otherwise all your requests could be denied.

### Dependences

go 1.2

github.com/bitly/go-simplejson

### Configuration

Look for `http.ListenAndServe(":19999", nil)`

### Sample

Input:

    curl http://127.0.0.1:19999 -d '{"wifi\_towers":[{"mac\_address":"00:27:22:2c:f1:bc"},{"mac\_address":"00:27:22:2c:f1:b5"},{"mac\_address":"00:27:22:2c:f1:eb"}],
    "cell\_towers":
    [{"mobile\_network\_code":"00","location\_area\_code":30571,"mobile\_country\_code":"460","cell\_id":46921},
    {"mobile\_network\_code":"00","location\_area\_code":30571,"mobile\_country\_code":"460","cell\_id":65535},
    {"mobile\_network\_code":"00","location\_area\_code":30571,"mobile\_country\_code":"460","cell\_id":65535},
    {"mobile\_network\_code":"00","location\_area\_code":30571,"mobile\_country\_code":"460","cell\_id":65535}]
    ,"radio\_type":"gsm","request\_address":false,"is\_chinacdma":0,"version":"1.1.0"}'

Output:

    {"access\_token":"","location":{"accuracy":1500,"altitude":0.0,"altitude\_accuracy":1500,"latitude":24.74439,"longitude":110.47237}}