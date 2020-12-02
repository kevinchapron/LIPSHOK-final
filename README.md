
# Framework Smart-Home Kit (FSHK)

[![Licence Apache2](https://img.shields.io/hexpm/l/plug.svg)](http://www.apache.org/licenses/LICENSE-2.0)

![Logo Golang][logo_golang]

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

The second system is the algorithm PRESENT[^1], which is an ultra lightweight algorithm, created to be run on low computational units. 
In some of our cases, such as BLE communications, remote devices have not the memory not the power to compute AES-128 fast enough. 
In those cases, as the data have very strict limitations, it is an ideal pick.

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


[logo_golang]:data:image/jpeg;base64,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
[^1]:Bogdanov A. et al. (2007) PRESENT: An Ultra-Lightweight Block Cipher. In: Paillier P., Verbauwhede I. (eds) Cryptographic Hardware and Embedded Systems - CHES 2007. CHES 2007. Lecture Notes in Computer Science, vol 4727. Springer, Berlin, Heidelberg. 
