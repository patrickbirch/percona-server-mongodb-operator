.. rn:: 1.0.0

Percona Kubernetes Operator for Percona Server for MongoDB
===========================================================

Percona announces the general availability of |Percona Kubernetes Operator for Percona Server for MongoDB| 1.0.0 on May 24, 2019. This release is now the current GA release in the 1.0 series. Download the latest version from the Percona Software Repositories. Please see the `GA release announcement`. All of Percona's software is open-source and free.

Installation
------------

Installation is performed by accessing the `Percona Software Repositories <https://www.percona.com/doc/kubernetes-operator-for-mongodb/kubernetes.html>`__ for Kubernetes and `OpenShift <https://www.percona.com/doc/kubernetes-operator-for-mongodb/openshift.html`__.

Notable Issues in Features
--------------------------
* :psbug:`237 <https://jira.percona.com/browse/CLOUD-237>`__ Use SSL/TLS connections for communications within the cluster.


 * :psbug:`219 <https://jira.percona.com/browse/CLOUD-219>`__ Automatic backup and restore provides a method of performing a hot backup of your MongoDB data while the system is running. `Percona XtraBackup <https://www.percona.com/software/mysql-database/percona-xtrabackup>`__ is a free, online, open-source, backup tool.


 * :psbug:`178 < https://jira.percona.com/browse/CLOUD-178>`__ Allows a `manual certificate signing request <https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/#create-a-certificate-signing-request-object-to-send-to-the-kubernetes-api>`_ to send to the Kubernetes API.
