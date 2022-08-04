/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	databasev1alpha1 "github.com/mmontes11/mariadb-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("BackupMariaDB controller", func() {
	Context("When creating a BackupMariaDB", func() {
		It("Should reconcile", func() {
			By("Creating BackupMariaDB")
			backupKey := types.NamespacedName{
				Name:      "backup-test",
				Namespace: defaultNamespace,
			}
			backup := databasev1alpha1.BackupMariaDB{
				ObjectMeta: metav1.ObjectMeta{
					Name:      backupKey.Name,
					Namespace: backupKey.Namespace,
				},
				Spec: databasev1alpha1.BackupMariaDBSpec{
					MariaDBRef: corev1.LocalObjectReference{
						Name: mariaDbName,
					},
					Storage: databasev1alpha1.Storage{
						ClassName: defaultStorageClass,
						Size:      storageSize,
					},
				},
			}
			Expect(k8sClient.Create(ctx, &backup)).To(Succeed())

			By("Expecting to create a Job eventually")
			Eventually(func() bool {
				var job batchv1.Job
				if err := k8sClient.Get(ctx, backupKey, &job); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Expecting BackupMariaDB to be complete eventually")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, backupKey, &backup); err != nil {
					return false
				}
				return backup.IsComplete()
			}, timeout, interval).Should(BeTrue())

			By("Deleting BackupMariaDB")
			Expect(k8sClient.Get(ctx, backupKey, &backup)).To(Succeed())
		})
	})
})