package engine

import (
	"github.com/go-logr/logr"
	autogenv1 "github.com/kyverno/kyverno/pkg/autogen/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/internal"
)

// GenerateResponse checks for validity of generate rule on the resource
func (e *engine) generateResponse(
	logger logr.Logger,
	policyContext engineapi.PolicyContext,
) engineapi.PolicyResponse {
	resp := engineapi.NewPolicyResponse()
	for _, rule := range autogenv1.ComputeRules(policyContext.Policy(), "") {
		logger := internal.LoggerWithRule(logger, rule)
		if ruleResp := e.filterRule(rule, logger, policyContext); ruleResp != nil {
			resp.Rules = append(resp.Rules, *ruleResp)
		}
	}
	return resp
}
