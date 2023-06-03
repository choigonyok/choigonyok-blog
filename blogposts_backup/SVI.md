TAG : network, study

(이미지1)

네트워크를 공부하다가 이상한 걸 발견
L2 스위치에 IP addr를 왜 할당? L2 스위치는 MAC address 기반이라 IP주소를 할당하지 못하는 걸로 알고있었는데?

찾아보니
SNMP, ssh, telnet등으로 원격관리를 하기 위해서 필요하다라고 함

근데 아~무도 왜 필요한지, IP address 할당과 SNMP, SSH가 무슨 상관인지는 아무도 안알려줌. 혼란. 왜??

SNMP, SSH, TELNET은 모두 IP address를 기반으로 통신함.
돌이켜보면, 블로그 배포한 aws ec2 ubuntu에 접속할 때도, IP address를 통해서 접속 했었음.
그래서 원격으로 관리하려면 IP를 할당해야하는구나~

기본적으로 스위치는 L2장치.
L3 스위치 (= multilayer switch)도 설정를 해주기 전까진 기본적으로 L2로 동작
한 스위치 안에 설정된 각기 다른 vlan끼리 통신하려면? 물리적으로는 연결되어있지만 논리적으로 다른 네트워크이기때문에 라우터가 필요. 그래서 스위치에 가상의 인터페이스를 만들어서 서로 다른 vlan끼리 라우팅을 하게 해주는 것이 바로 SVI(switch virtual interface)
(https://www.networkstraining.com/what-is-cisco-svi-configuration-example/)

이 SVI를 L3가 아닌 L2에서 설정하면, L2 스위치의 SVI는 L3 이기 때문에 IP주소를 갖고있고, 이 IP주소를 통해서 L2스위치에 원격으로 접속해서 L2 스위치를 원격관리할 수 있게 되는 것!

l3 스위치는 l2역할을 겸할 수 있음
L2에서 SVI는 원격관리를 위해서
L3에서 SVI는 원격관리 + 라우팅(다른 네트워크와의 연결) 위해서

모든 L2 스위치가 SVI가 된다고 다 라우팅이 가능한 L3스위치가 되는 것은 아님
1. L2에서 SVI가 설정되면 라우팅이 가능한 스위치 = L3 스위치(Multilayer 스위치), SVI는 원격관리+라우팅 
2. L2에서 SVI가 설정되어도 라우팅이 불가능한 스위치 = SVI는 원격관리만
스위치의 모델과 기능에 따라 차이가 있는 듯 하다.

### SVI로 스위치에 IP를 할당하면 생기는 장점

라우터를 거치는 것보다 훨씬 빠름
https://ipwithease.com/switched-virtual-interface-svi/