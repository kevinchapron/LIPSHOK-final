
# Framework Smart-Home Kit (FSHK)

[![Licence Apache2](https://img.shields.io/hexpm/l/plug.svg)](http://www.apache.org/licenses/LICENSE-2.0)

![Logo Golang](https://github.com/golang/go/blob/master/doc/gopher/doc.png)

This thesis has been developed using GoLang.

---

How it works
------------

The framework is in the `main-app` folder, to manage the main server and its interaction. 

Once started, it will create its own database file using sqlite, with the name `FSHK.db`, wherever the program is started.
The program will put default values in it, but you can override them when it is shut down. On the restart of the program, it will take into account every data in it.

For example, using the table `database_device_types`, you can specify which protocol you want activated on the start of the service. If one of your device runs on this protocol, but you deactivated it, it won't be able to connect.
See independant categories for more explanation about each protocol management.

Messaging
---------

Each message must respect the same architecture, except for specified protocols. 
It must have at least 24 bytes for the message header, and a maximum of 65 536 bytes at most. So, body can contain up to 65 512 bytes.

Repartitions of the header : 
* **0-1  :** Data Length
* **2-14 :** Aes IV
* **15   :** Data Type
* **16-24:** Padding (reservation for possible upgrade) 

Of course the data length must match the data, otherwise the message won't be processed.

Security
--------

The system uses two distincts encryptions.

The first one, and the most used is AES-128. It is used for most of the communications of the kit, since most of it transits through big packets such as the ones in WebSockets (maxed at 65536 here). 
A new IV is generated every time for every message, while the MasterKey is hard-coded for every device. Since the encrypted data is not necessarily the same size of its original data, the first bytes must be edited to fit.

The second system is the algorithm PRESENT<a href="#note1" id="note1ref"><sup>1</sup></a>, which is an ultra lightweight algorithm, created to be run on low computational units. 
In some of our cases, such as BLE communications, remote devices have not the memory not the power to compute AES-128 fast enough. 
In those cases, as the data have very strict limitations, it is an ideal pick.

Protocols
---------

<details>
<summary><b>UDP</b></summary>

<p>
This software allows UDP packets to be received on the port <b>5010</b>
Then, it forwards it to the main app, to register everything.


If something is received, the service will acknowledge it using a return message <b>{"data":"OK"}</b> everytime. 
If you do not receive it, your message has not been received.
</p>
    
</details>

<details>
<summary><b>TCP</b></summary>

<p>
This software allows TCP packets to be received on the port <b>5020</b>
Then, it forwards it to the main app, to register everything.


If something is received, the service will acknowledge it using a return message <b>{"data":"OK"}</b> everytime. 
If you do not receive it, your message has not been received.
</p>
    
</details>



Web Interface
-------------

A web interface is available once the program is started, on the port 5002. It will interact with the server using a specific websocket designed to this effect. 


Authors
---
**[Kévin Chapron](http://kevin-chapron.fr/)**

License
---
    Copyright 2020 Kévin Chapron.

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

<a id="note1" href="#note1ref"><sup>1</sup></a> Bogdanov A. et al. (2007) PRESENT: An Ultra-Lightweight Block Cipher. In: Paillier P., Verbauwhede I. (eds) Cryptographic Hardware and Embedded Systems - CHES 2007. CHES 2007. Lecture Notes in Computer Science, vol 4727. Springer, Berlin, Heidelberg.
