package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

const NamespaceIndexName = "namespace"
const NodeNameIndexName = "nodeName"

func NamespaceIndexFunc(obj interface{}) ([]string, error) {
	m, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	return []string{m.GetNamespace()}, nil
}

func NodeNameIndexFunc(obj interface{}) ([]string, error) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return []string{""}, fmt.Errorf("obj is not pod")
	}
	return []string{pod.Spec.NodeName}, nil
}

func main() {
	index := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		NamespaceIndexName: NamespaceIndexFunc,
		NodeNameIndexName:  NodeNameIndexFunc,
	})

	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "index-pod-1",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			NodeName: "node1",
		},
	}
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "index-pod-2",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			NodeName: "node2",
		},
	}
	pod3 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "index-pod-3",
			Namespace: "kube-system",
		},
		Spec: v1.PodSpec{
			NodeName: "node2",
		},
	}

	_ = index.Add(pod1)
	_ = index.Add(pod2)
	_ = index.Add(pod3)

	pods, err := index.ByIndex(NamespaceIndexName, "default")
	if err != nil {
		fmt.Println(err)
	}

	for _, pod := range pods {
		fmt.Printf("%+v\n", pod)
	}

	pods, err = index.ByIndex(NodeNameIndexName, "node2")
	if err != nil {
		fmt.Println(err)
	}

	for _, pod := range pods {
		fmt.Printf("%+v\n", pod)
	}
}

//# Indexers 就是包含的所有索引器(分类)以及对应实现
//indexers: {
//	"namespace": NamespaceIndexFunc,
//	"nodeName": NodeNameIndexFunc,
//}
//
//# Indices 就是包含的所有索引分类中所有的索引数据
//indices: {
//	"namespace": {  #namespace 这个索引分类下的所有索引数据
//		"default": ["pod-1", "pod-2"],  # Index 就是一个索引键下所有的对象键列表
//		"kube-system": ["pod-3"]   # Index
//	},
//	"nodeName": {  # nodeName 这个索引分类下的所有索引数据(对象键列表)
//		"node1": ["pod-1"],  # Index
//		"node2": ["pod-2", "pod-3"]  # Index
//	}
//}
