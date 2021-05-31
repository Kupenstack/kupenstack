/*
Copyright 2021 The Kupenstack Authors.

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

// Package VM implements virtual-machine-reconciler for kupenstack controller.
//
// Events
//
// The following events are thrown by reconciler:
//  REASON                MESSAGE
//
//  ImageNotFound         Image %s not found.
//  KeyPairNotFound       Key pair %s not found.
//  FlavorNotFound        Flavor %s not found.
//  Created               Virtual Machine created.
//  CreateFailed          Virtual Machine create failed. error: %s
//  DeleteFailed          Vitual Machine deletion failed. error: %s
//  Deleted               Virtual Machine deleted.
package vm
