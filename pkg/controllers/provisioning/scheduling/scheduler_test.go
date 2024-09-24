package scheduling

import (
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/karpenter/pkg/apis/v1beta1"
	"sigs.k8s.io/karpenter/pkg/scheduling"
)

func TestGetDaemonOverhead(t *testing.T) {
	nodeClaimTemplates := []*NodeClaimTemplate{
		{
			NodeClaimTemplate: v1beta1.NodeClaimTemplate{
				Spec: v1beta1.NodeClaimSpec{
					NodeClassRef: &v1beta1.NodeClassReference{
						Name:       "test",
						Kind:       "test",
						APIVersion: "group/test",
					},
					Taints: []v1.Taint{
						{
							Key:    "example.com/no-schedule",
							Value:  "true",
							Effect: v1.TaintEffectPreferNoSchedule,
						},
					},
				},
			},
			Requirements: scheduling.Requirements{},
		},
		{
			NodeClaimTemplate: v1beta1.NodeClaimTemplate{
				Spec: v1beta1.NodeClaimSpec{
					NodeClassRef: &v1beta1.NodeClassReference{
						Name:       "test2",
						Kind:       "test2",
						APIVersion: "group/test2",
					},
					Taints: []v1.Taint{
						{
							Key:    "example.com/no-schedule",
							Value:  "true",
							Effect: v1.TaintEffectPreferNoSchedule,
						},
					},
				},
			},
			Requirements: scheduling.Requirements{},
		},
	}
	nodeClaimTemplates[0].Requirements.Add(scheduling.NewRequirement("role", v1.NodeSelectorOpIn, "monitor"))
	nodeClaimTemplates[0].Requirements.Add(scheduling.NewRequirement("owner", v1.NodeSelectorOpIn, "kcc"))
	nodeClaimTemplates[1].Requirements.Add(scheduling.NewRequirement("consumer", v1.NodeSelectorOpIn, "ngcweb"))

	daemonSetPods := []*v1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "daemon-pod-1"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "example-container",
						Image: "nginx:latest",
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse("500m"),
								v1.ResourceMemory: resource.MustParse("256Mi"),
							},
						},
					},
				},
				Tolerations: []v1.Toleration{
					{Key: "example.com/no-schedule", Operator: v1.TolerationOpExists, Effect: v1.TaintEffectPreferNoSchedule},
				},
				Affinity: &v1.Affinity{
					NodeAffinity: &v1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
							NodeSelectorTerms: []v1.NodeSelectorTerm{
								{
									MatchExpressions: []v1.NodeSelectorRequirement{
										{
											Key:      "os",
											Operator: v1.NodeSelectorOpIn,
											Values:   []string{"linux"},
										},
									},
								},
								{
									MatchExpressions: []v1.NodeSelectorRequirement{
										{
											Key:      "role",
											Operator: v1.NodeSelectorOpIn,
											Values:   []string{"monitor"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "daemon-pod-2"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "example-container-2",
						Image: "nginx:latest",
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse("100m"),
								v1.ResourceMemory: resource.MustParse("32Mi"),
							},
						},
					},
				},
				Tolerations: []v1.Toleration{
					{Key: "example.com/no-schedule", Operator: v1.TolerationOpExists, Effect: v1.TaintEffectPreferNoSchedule},
				},
				Affinity: &v1.Affinity{
					NodeAffinity: &v1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
							NodeSelectorTerms: []v1.NodeSelectorTerm{
								{
									MatchExpressions: []v1.NodeSelectorRequirement{
										{
											Key:      "role",
											Operator: v1.NodeSelectorOpIn,
											Values:   []string{"monitor"},
										},
									},
								},
								{
									MatchExpressions: []v1.NodeSelectorRequirement{
										{
											Key:      "owner",
											Operator: v1.NodeSelectorOpIn,
											Values:   []string{"kcc"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "daemon-pod-3"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "example-container-3",
						Image: "nginx:latest",
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse("200m"),
								v1.ResourceMemory: resource.MustParse("64Mi"),
							},
						},
					},
				},
				Tolerations: []v1.Toleration{
					{Key: "example.com/no-schedule", Operator: v1.TolerationOpExists, Effect: v1.TaintEffectPreferNoSchedule},
				},
				Affinity: &v1.Affinity{
					NodeAffinity: &v1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
							NodeSelectorTerms: []v1.NodeSelectorTerm{
								{
									MatchExpressions: []v1.NodeSelectorRequirement{
										{
											Key:      "os",
											Operator: v1.NodeSelectorOpIn,
											Values:   []string{"linux"},
										},
									},
								},
								{
									MatchExpressions: []v1.NodeSelectorRequirement{
										{
											Key:      "owner",
											Operator: v1.NodeSelectorOpIn,
											Values:   []string{"kcs"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	overhead := getDaemonOverhead(nodeClaimTemplates, daemonSetPods)
	fmt.Println("Map of nodeClaimTemplates:\n", overhead)

}
