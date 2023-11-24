# AsciiArtify Kubernetes Deployment Tools Evaluation

## Introduction
AsciiArtify, a startup focused on developing a new software product for transforming images into ASCII art using Machine Learning, faces the challenge of selecting the right tool for local Kubernetes cluster development. The team, comprised of two young programmers with expertise in software development but lacking DevOps experience, is considering three options: minikube, kind and k3d.

## Characteristics
### Minikube
- `Supported OS and Architectures:` Works on multiple operating systems, including Windows, macOS, and Linux. Supports various architectures.  
- `Automation Capability:` Offers automation for cluster creation and management.
- `Additional Features:` Suitable for local development and testing. Concerns arise regarding scalability limitations.  

### Kind (Kubernetes IN Docker)
- `Supported OS and Architectures:` Compatible with Windows, macOS, and Linux. Works within Docker containers.  
- `Automation Capability:` Allows the creation of local Kubernetes clusters in Docker containers.  
- `Additional Features:` Considered for local testing purposes.

### k3d
- `Supported OS and Architectures:` Works on multiple operating systems. Uses Rancher Kubernetes Engine (RKE) in Docker containers.
- `Automation Capability:` Facilitates quick creation and testing of Kubernetes clusters in Docker containers.
- `Additional Features:` Chosen for preparing Proof of Concept (PoC).

### Characteristic

| **Pros and Cons**                               | **Minikube**                                     | **Kind**                                         | **k3d**                                          | **Podman**                                       |
|--------------------------------------------------|--------------------------------------------------|--------------------------------------------------|--------------------------------------------------|--------------------------------------------------|
| **Pros**                                      | + Easy to use<br>+ Suitable for local development and testing | + Suitable for local development and testing<br>+ Works within Docker containers<br>+ Possibility for local testing | + Suitable for local development and testing<br>+ Works within Docker containers<br>+ Fast cluster creation and testing | + Easy to use<br>+ Suitable for local development and testing<br>+ Works within Docker containers<br>+ Light alternative to Docker |
| **Cons**                                      | - Doubts about scalability<br>- Potential limitations | - Limited information on scalability<br>- Limited community documentation | - Limited documentation<br>- Potential scalability concerns | - Limited information on scalability<br>- Limited community documentation |


## Demonstration
Recommended Tool: k3d  Deployment of "Hello World" Application on Kubernetes  

![Application on Kubernetes](622883.gif)  


## Conclusion
After practical exploration, k3d stands out as the recommended tool for AsciiArtify's PoC. Its quick cluster creation and simplicity make it suitable for initial development. However, it's crucial to consider the limited community documentation and potential scalability concerns. Additionally, Podman is introduced as a lightweight alternative to Docker, offering rootless containers and direct integration with systemd, although with a less mature ecosystem. AsciiArtify should carefully weigh the pros and cons before making a final decision.