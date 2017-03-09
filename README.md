# Podder
Podder is a rudimentary PaaS that can directly deploy binaries in a cluster and orchestrate them. 


> As computation continues to move into the cloud, the computing platform of interest no longer resembles a pizza box or a refrigerator, but a warehouse full of computers. These new large datacenters are quite different from traditional hosting facilities of earlier times and cannot be viewed simply as a collection of co-located servers. Large portions of the hardware and software resources in these facilities must work in concert to efficiently deliver good levels of Internet service performance, something that can only be achieved by a holistic approach to their design and deployment. In other words, we must treat the datacenter itself as one massive warehouse-scale computer (WSC) -- Luiz André Barroso, Urs Hölzle

Docker marked an epoch in the way we did deployments and gave us a new primitive which we can schedule across multiple nodes without worrying about a lot of things that troubled distributed operating systems that preceded Docker. Many systems strive to acheive the goal of treating the entire datacenter as a computer. Google's Borg has been running in production for so many years. Kubernetes uses Docker containers as a primitive and treats the entire datacenter as a single computer. Kubernetes is often called the kernel of the datacenter. It's not far-feteched as it does almost everything an operating system kernel is supposed to do. Kubernetes manages resources and schedules compute across nodes, much like what a normal operating system would. 

Podder is built on top of Kubernetes. Podder is an attempt at being the shell that one uses on top of the datacenter's kernel. Podder can reduce the barrier to start using advancements in distributed computing by giving the masses a simple interface to deploy their binaries. 

#How to use
With Podder running on your kubernetes cluster, if you wanted to deploy a binary, you simply drag-drop the binary into the Podder interface and it automatically creates the replication controller and a service. This is extremely handy for someone trying to build simple web applications or microservices. 



#Todo

- [ ] Setup Ingress controller
- [ ] Setup Stateful sets for database and object store
- [ ] Update README.