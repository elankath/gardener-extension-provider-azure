// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"fmt"

	"github.com/gardener/gardener/extensions/pkg/controller"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	api "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure"
	"github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/install"
	apiv1alpha1 "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
)

var (
	// Scheme is a scheme with the types relevant for Azure actuators.
	Scheme *runtime.Scheme

	decoder runtime.Decoder

	// lenientDecoder is a decoder that does not use strict mode.
	lenientDecoder runtime.Decoder

	// InfrastructureStateTypeMeta is the TypeMeta of the Azure InfrastructureStatus
	InfrastructureStateTypeMeta = metav1.TypeMeta{
		APIVersion: apiv1alpha1.SchemeGroupVersion.String(),
		Kind:       "InfrastructureState",
	}
)

func init() {
	Scheme = runtime.NewScheme()
	utilruntime.Must(install.AddToScheme(Scheme))

	decoder = serializer.NewCodecFactory(Scheme, serializer.EnableStrict).UniversalDecoder()
	lenientDecoder = serializer.NewCodecFactory(Scheme).UniversalDecoder()
}

// InfrastructureConfigFromInfrastructure extracts the InfrastructureConfig from the
// ProviderConfig section of the given Infrastructure.
func InfrastructureConfigFromInfrastructure(infra *extensionsv1alpha1.Infrastructure) (*api.InfrastructureConfig, error) {
	config := &api.InfrastructureConfig{}
	if infra.Spec.ProviderConfig != nil && infra.Spec.ProviderConfig.Raw != nil {
		if _, _, err := decoder.Decode(infra.Spec.ProviderConfig.Raw, nil, config); err != nil {
			return nil, err
		}
		return config, nil
	}
	return nil, fmt.Errorf("provider config is not set on the infrastructure resource")
}

// InfrastructureStatusFromRaw extracts the InfrastructureStatus from the
// ProviderStatus section of the given Infrastructure.
func InfrastructureStatusFromRaw(raw *runtime.RawExtension) (*api.InfrastructureStatus, error) {
	config := &api.InfrastructureStatus{}
	if raw != nil && raw.Raw != nil {
		if _, _, err := lenientDecoder.Decode(raw.Raw, nil, config); err != nil {
			return nil, err
		}
		return config, nil
	}
	return nil, fmt.Errorf("provider status is not set on the infrastructure resource")
}

// CloudProfileConfigFromCluster decodes the provider specific cloud profile configuration for a cluster
func CloudProfileConfigFromCluster(cluster *controller.Cluster) (*api.CloudProfileConfig, error) {
	var cloudProfileConfig *api.CloudProfileConfig
	if cluster != nil && cluster.CloudProfile != nil && cluster.CloudProfile.Spec.ProviderConfig != nil && cluster.CloudProfile.Spec.ProviderConfig.Raw != nil {
		cloudProfileConfig = &api.CloudProfileConfig{}
		if _, _, err := decoder.Decode(cluster.CloudProfile.Spec.ProviderConfig.Raw, nil, cloudProfileConfig); err != nil {
			return nil, fmt.Errorf("could not decode providerConfig of cloudProfile for '%s': %w", kutil.ObjectName(cluster.CloudProfile), err)
		}
	}
	return cloudProfileConfig, nil
}

// BackupConfigFromBackupBucket decodes the provider specific config from a given BackupBucket object.
func BackupConfigFromBackupBucket(backupBucket *extensionsv1alpha1.BackupBucket) (*api.BackupConfig, error) {
	var backupConfig *api.BackupConfig
	if backupBucket != nil && backupBucket.Spec.DefaultSpec.ProviderConfig != nil {
		backupConfig = &api.BackupConfig{}
		if _, _, err := decoder.Decode(backupBucket.Spec.ProviderConfig.Raw, nil, backupConfig); err != nil {
			return nil, err
		}
	}
	return backupConfig, nil
}

// BackupConfigFromBackupEntry  decodes the provider specific config from a given BackupEntry object.
func BackupConfigFromBackupEntry(backupEntry *extensionsv1alpha1.BackupEntry) (*api.BackupConfig, error) {
	var backupConfig *api.BackupConfig
	if backupEntry != nil && backupEntry.Spec.DefaultSpec.ProviderConfig != nil {
		backupConfig = &api.BackupConfig{}
		if _, _, err := decoder.Decode(backupEntry.Spec.ProviderConfig.Raw, nil, backupConfig); err != nil {
			return nil, err
		}
	}
	return backupConfig, nil
}

// InfrastructureStateFromRaw extracts the state from the Infrastructure. If no state was available, it returns a "zero" value InfrastructureState object.
func InfrastructureStateFromRaw(raw *runtime.RawExtension) (*api.InfrastructureState, error) {
	state := &api.InfrastructureState{}
	if raw != nil && raw.Raw != nil {
		if _, _, err := lenientDecoder.Decode(raw.Raw, nil, state); err != nil {
			return nil, err
		}
	}

	if state.Data == nil {
		state.Data = make(map[string]string)
	}
	if state.ManagedItems == nil {
		state.ManagedItems = make([]api.AzureResource, 0)
	}

	return state, nil
}
