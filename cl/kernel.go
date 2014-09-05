/*
 * Copyright © 2012 go-opencl authors
 *
 * This file is part of go-opencl.
 *
 * go-opencl is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * go-opencl is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with go-opencl.  If not, see <http://www.gnu.org/licenses/>.
 */

package cl

/*
#cgo CFLAGS: -I CL
#cgo linux LDFLAGS: -lOpenCL
#cgo windows LDFLAGS: -lOpenCL
#cgo darwin LDFLAGS: -framework OpenCL

#ifdef MAC
	#include "OpenCL/cl.h"
#else
	#include "CL/opencl.h"
#endif //MAC

*/
import "C"

import (
	"unsafe"
)

type KernelProperty C.cl_kernel_info

const (
	KERNEL_FUNCTION_NAME   KernelProperty = C.CL_KERNEL_FUNCTION_NAME
	KERNEL_NUM_ARGS        KernelProperty = C.CL_KERNEL_NUM_ARGS
	KERNEL_REFERENCE_COUNT KernelProperty = C.CL_KERNEL_REFERENCE_COUNT
	KERNEL_CONTEXT         KernelProperty = C.CL_KERNEL_CONTEXT
	KERNEL_PROGRAM         KernelProperty = C.CL_KERNEL_PROGRAM
	// new in 1.2
	// KERNEL_ATTRIBUTES      KernelProperty = C.CL_KERNEL_ATTRIBUTES
)

func (k *Kernel) Property(prop KernelProperty) interface{} {
	if value, ok := k.properties[prop]; ok {
		return value
	}

	var data interface{}
	var length C.size_t
	var ret C.cl_int

	switch prop {
	case
		KERNEL_FUNCTION_NAME:
		if ret = C.clGetKernelInfo(k.id, C.cl_kernel_info(prop), 0, nil, &length); ret != C.CL_SUCCESS || length < 1 {
			data = ""
			break
		}

		buf := make([]C.char, length)
		if ret = C.clGetKernelInfo(k.id, C.cl_kernel_info(prop), length, unsafe.Pointer(&buf[0]), &length); ret != C.CL_SUCCESS || length < 1 {
			data = ""
			break
		}
		data = C.GoStringN(&buf[0], C.int(length-1))

	case
		KERNEL_NUM_ARGS,
		KERNEL_REFERENCE_COUNT:
		var val C.cl_uint
		ret = C.clGetKernelInfo(k.id, C.cl_kernel_info(prop), C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), &length)
		data = val

	case
		KERNEL_CONTEXT:
		var val C.cl_context
		ret = C.clGetKernelInfo(k.id, C.cl_kernel_info(prop), C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), &length)
		data = val

	case
		KERNEL_PROGRAM:
		var val C.cl_program
		ret = C.clGetKernelInfo(k.id, C.cl_kernel_info(prop), C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), &length)
		data = val

	// new in 1.2
	// case
	// 	KERNEL_ATTRIBUTES:
	// 	if ret = C.clGetKernelInfo(k.id, C.cl_kernel_info(prop), 0, nil, &length); ret != C.CL_SUCCESS || length < 1 {
	// 		data = ""
	// 		break
	// 	}

	// 	buf := make([]C.char, length)
	// 	if ret = C.clGetKernelInfo(k.id, C.cl_kernel_info(prop), length, unsafe.Pointer(&buf[0]), &length); ret != C.CL_SUCCESS || length < 1 {
	// 		data = ""
	// 		break
	// 	}
	// 	data = C.GoStringN(&buf[0], C.int(length-1))

	default:
		return nil
	}

	if ret != C.CL_SUCCESS {
		return nil
	}
	k.properties[prop] = data
	return k.properties[prop]
}

type KernelWorkGroupProperty C.cl_kernel_work_group_info

const (
	KERNEL_WORK_GROUP_SIZE             KernelWorkGroupProperty = C.CL_KERNEL_WORK_GROUP_SIZE
	COMPILE_WORK_GROUP_SIZE            KernelWorkGroupProperty = C.CL_KERNEL_COMPILE_WORK_GROUP_SIZE
	LOCAL_MEM_SIZE                     KernelWorkGroupProperty = C.CL_KERNEL_LOCAL_MEM_SIZE
	PREFERRED_WORK_GROUP_SIZE_MULTIPLE KernelWorkGroupProperty = C.CL_KERNEL_PREFERRED_WORK_GROUP_SIZE_MULTIPLE
	PRIVATE_MEM_SIZE                   KernelWorkGroupProperty = C.CL_KERNEL_PRIVATE_MEM_SIZE
)

func (k *Kernel) WorkGroupProperty(d Device, prop KernelWorkGroupProperty) interface{} {
	if value, ok := k.workgroupproperties[prop]; ok {
		return value
	}

	var data interface{}
	var length C.size_t
	var ret C.cl_int

	switch prop {
	case KERNEL_WORK_GROUP_SIZE,
		PREFERRED_WORK_GROUP_SIZE_MULTIPLE:
		var val C.size_t
		ret = C.clGetKernelWorkGroupInfo(k.id, d.id, C.cl_kernel_work_group_info(prop), C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), &length)
		data = val

	case COMPILE_WORK_GROUP_SIZE:
		// size_t[3]
		break

	case LOCAL_MEM_SIZE,
		PRIVATE_MEM_SIZE:
		var val C.cl_ulong
		ret = C.clGetKernelWorkGroupInfo(k.id, d.id, C.cl_kernel_work_group_info(prop), C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), &length)
		data = val

	default:
		return nil
	}

	if ret != C.CL_SUCCESS {
		return nil
	}
	k.workgroupproperties[prop] = data
	return k.workgroupproperties[prop]
}

type Kernel struct {
	id                  C.cl_kernel
	properties          map[KernelProperty]interface{}
	workgroupproperties map[KernelWorkGroupProperty]interface{}
}

func (k *Kernel) SetArg(index uint, arg interface{}) error {
	var ret C.cl_int

	switch t := arg.(type) {
	case *Buffer:
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(t.id)), unsafe.Pointer(&t.id))
	case *Image:
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(t.id)), unsafe.Pointer(&t.id))
	case *Sampler:
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(t.id)), unsafe.Pointer(&t.id))
	case int32:
		f := C.int(t)
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(f)), unsafe.Pointer(&f))
	case float32:
		f := C.float(t)
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(f)), unsafe.Pointer(&f))

	default:
		return Cl_error(C.CL_INVALID_VALUE)
	}

	if ret != C.CL_SUCCESS {
		return Cl_error(ret)
	}
	return nil
}

func (k *Kernel) SetArgs(offset uint, args []interface{}) error {
	for i, arg := range args {
		if err := k.SetArg(offset+uint(i), arg); err != nil {
			return err
		}
	}
	return nil
}

func (k *Kernel) release() error {
	if k.id != nil {
		if err := C.clReleaseKernel(k.id); err != C.CL_SUCCESS {
			return Cl_error(err)
		}
		k.id = nil
	}
	return nil
}
