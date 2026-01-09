# 🔄 SDD Workflow Reference

The Ultimate SDD Framework enforces a structured development workflow with mandatory phases and approval gates. This document provides detailed guidance for each phase of the Spec-Driven Development process.

## Workflow Overview

```
🎯 Specify → 🏗️ Plan → ✅ Approve → 📋 Task → 💻 Execute → 🔍 Review → ✅ Complete
```

Each phase has specific inputs, outputs, and validation requirements.

## Phase 1: Specify 🎯

**Command:** `sdd specify "feature description"`

**Agent:** Product Manager (`pm.md`)

**Purpose:** Convert high-level ideas into detailed technical specifications

### Inputs
- User feature request or description
- Existing codebase context (via LSP)
- Current project requirements

### Activities
1. **Requirements Analysis**
   - Identify functional requirements
   - Define non-functional requirements
   - Document business rules and constraints

2. **Edge Case Identification**
   - Consider error scenarios
   - Identify boundary conditions
   - Document exceptional flows

3. **Acceptance Criteria Definition**
   - Write testable success conditions
   - Define done criteria
   - Establish validation methods

### Outputs
- `spec.md`: Detailed feature specification with:
  - Overview and objectives
  - Functional requirements
  - Non-functional requirements
  - Constraints and limitations
  - Acceptance criteria
  - Edge cases and risk assessment

### Validation
- Specification covers all aspects of the request
- Requirements are clear and testable
- Edge cases are identified
- Acceptance criteria are measurable

---

## Phase 2: Plan 🏗️

**Command:** `sdd plan`

**Agent:** System Architect (`architect.md`)

**Purpose:** Design the technical architecture and implementation approach

### Inputs
- Approved feature specification (`spec.md`)
- Existing codebase analysis (via LSP)
- Technology stack and constraints

### Activities
1. **Architecture Design**
   - Define system components and boundaries
   - Select technology stack and frameworks
   - Design data architecture and storage

2. **Integration Planning**
   - Identify external system dependencies
   - Plan API designs and contracts
   - Define data flow and processing

3. **Risk Assessment**
   - Identify technical risks and challenges
   - Define mitigation strategies
   - Estimate implementation complexity

### Outputs
- `plan.md`: Comprehensive architecture plan with:
  - System overview and component diagram
  - Technology stack rationale
  - Data architecture and flow
  - API design specifications
  - Implementation phases and priorities
  - Risk assessment and mitigation

### Validation
- Architecture addresses all requirements
- Technology choices are justified
- Integration points are clearly defined
- Risks are identified and mitigated

---

## Phase 3: Approve ✅ (Gate)

**Command:** `sdd approve`

**Agent:** Human Developer or QA Agent

**Purpose:** Quality gate ensuring plans meet standards before implementation

### The "Merged Secret"
**You cannot proceed to task breakdown until the architecture plan is approved.**

This prevents poorly designed features from being implemented and ensures:
- Plans meet quality standards
- Requirements are properly addressed
- Technical approach is sound
- Risks are acceptable

### Approval Criteria
- [ ] Requirements fully addressed in plan
- [ ] Technology choices appropriate for project
- [ ] Architecture scalable and maintainable
- [ ] Security considerations included
- [ ] Performance requirements met
- [ ] Integration risks identified
- [ ] Implementation approach feasible

### Approval Methods
1. **Human Approval:** Developer reviews and approves
2. **Automated Checks:** Basic validation rules
3. **Peer Review:** Team member approval
4. **QA Review:** Quality assurance validation

---

## Phase 4: Task 📋

**Command:** `sdd task`

**Agent:** Software Developer (`developer.md`)

**Purpose:** Break down the approved plan into specific, actionable tasks

### Prerequisites
- ✅ Architecture plan approved (`sdd approve`)

### Inputs
- Approved architecture plan (`plan.md`)
- Implementation complexity assessment
- Team capacity and skills

### Activities
1. **Task Breakdown**
   - Decompose plan into manageable units
   - Estimate time and complexity
   - Identify dependencies and blockers

2. **Implementation Planning**
   - Define development sequence
   - Assign priorities and milestones
   - Plan testing and validation

3. **Quality Assurance Planning**
   - Define testing requirements
   - Plan code review processes
   - Establish quality checkpoints

### Outputs
- `tasks.md`: Detailed task breakdown with:
  - Task checklist with IDs and descriptions
  - Time estimates and priorities
  - Dependencies and prerequisites
  - Acceptance criteria per task
  - Testing requirements
  - Quality gates

### Task Template
```markdown
## Implementation Tasks

### 🔧 Infrastructure & Setup
- [ ] **TASK-001**: Set up project structure and dependencies
  - Create Go module structure
  - Initialize database schema
  - Set up Docker configuration
  - Configure CI/CD pipeline
  - **Acceptance Criteria**: Project builds successfully, tests pass
  - **Estimated Time**: 4 hours
  - **Priority**: High

### 🏗️ Core API Development
- [ ] **TASK-002**: Implement REST API endpoints
  - CRUD operations for main entities
  - Input validation and error handling
  - Response formatting
  - API documentation
  - **Acceptance Criteria**: All endpoints return correct responses
  - **Estimated Time**: 8 hours
  - **Priority**: High
```

### Validation
- Tasks cover complete implementation
- Estimates are realistic and detailed
- Dependencies are identified
- Acceptance criteria are testable

---

## Phase 5: Execute 💻

**Command:** `sdd execute`

**Agent:** Software Developer (`developer.md`)

**Purpose:** Provide detailed implementation guidance and development standards

### Inputs
- Task breakdown (`tasks.md`)
- Architecture plan (`plan.md`)
- Existing codebase patterns

### Activities
1. **Development Environment Setup**
   - Configure local development environment
   - Set up databases and external services
   - Install dependencies and tools

2. **Implementation Guidance**
   - Provide code structure and organization
   - Define coding standards and patterns
   - Guide testing approaches

3. **Quality Assurance Integration**
   - Define testing strategies
   - Set up code quality tools
   - Plan deployment and monitoring

### Outputs
- `implementation.md`: Comprehensive development guide with:
  - Environment setup instructions
  - Development workflow and standards
  - Code organization and patterns
  - Testing strategies and frameworks
  - Deployment and production considerations
  - Quality assurance procedures

### Development Standards
```markdown
### Code Quality Checklist
#### Before Commit
- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] Code follows style guidelines
- [ ] No TODO comments left unresolved
- [ ] Documentation updated

#### Code Review Requirements
- [ ] Logic is correct and efficient
- [ ] Error handling is appropriate
- [ ] Security considerations addressed
- [ ] Performance implications reviewed
- [ ] Tests are comprehensive
```

---

## Phase 6: Review 🔍

**Command:** `sdd review`

**Agent:** Quality Assurance (`qa.md`)

**Purpose:** Conduct comprehensive quality assessment before completion

### Inputs
- Implementation progress and artifacts
- Task completion status
- Test results and coverage
- Code quality metrics

### Activities
1. **Code Quality Assessment**
   - Review code structure and organization
   - Validate adherence to standards
   - Check for security vulnerabilities
   - Assess maintainability and readability

2. **Testing Validation**
   - Verify test coverage and quality
   - Validate testing strategies
   - Check integration and system tests
   - Assess performance benchmarks

3. **Requirements Compliance**
   - Verify implementation meets specifications
   - Validate acceptance criteria
   - Check edge cases and error handling
   - Assess user experience quality

### Outputs
- `review.md`: Quality assurance report with:
  - Executive summary and quality score
  - Detailed findings and issues
  - Recommendations and improvements
  - Approval recommendation
  - Deployment readiness assessment

### Quality Score Framework
```markdown
### Quality Score: A+ (96/100)

| Category | Score | Weight | Weighted |
|----------|-------|--------|----------|
| Requirements | 100%% | 25%% | 25 |
| Code Quality | 95%% | 20%% | 19 |
| Testing | 98%% | 20%% | 19.6 |
| Security | 100%% | 15%% | 15 |
| Performance | 92%% | 10%% | 9.2 |
| Documentation | 95%% | 10%% | 9.5 |
| **Total** | | | **96.3** |
```

### Validation
- All acceptance criteria met
- Code quality standards satisfied
- Testing coverage adequate
- Security requirements fulfilled
- Performance acceptable
- Documentation complete

---

## Phase 7: Complete ✅

**Command:** `sdd approve` (final approval)

**Agent:** Human Developer or QA Agent

**Purpose:** Final approval for production deployment

### Final Validation
- [ ] All tasks completed and tested
- [ ] QA review passed with acceptable score
- [ ] Documentation updated and accurate
- [ ] Deployment procedures documented
- [ ] Rollback plan prepared
- [ ] Monitoring and alerting configured

### Completion Activities
1. **Deployment Preparation**
   - Final testing in staging environment
   - Performance validation
   - Security scanning

2. **Documentation Updates**
   - Update API documentation
   - Update deployment guides
   - Document known issues and limitations

3. **Knowledge Transfer**
   - Team training and handoff
   - Operational procedures documented
   - Support and maintenance guidelines

---

## Workflow States & Transitions

### State Diagram
```
Init ──────→ Specify ──────→ Plan ──────→ Approve ──────→ Task ──────→ Execute ──────→ Review ──────→ Complete
   │           │             │             │              │             │              │              │
   │           │             │             │              │             │              │              │
   └───────────┼─────────────┼─────────────┼──────────────┼─────────────┼──────────────┼──────────────┘
               │             │             │              │             │              │
               └─────────────┼─────────────┼──────────────┼─────────────┼──────────────┘
                             │             │              │             │
                             └─────────────┼──────────────┼─────────────┼──────────────┘
                                           │              │             │
                                           └──────────────┼─────────────┼──────────────┘
                                                          │             │
                                                          └─────────────┼──────────────┘
                                                                        │
                                                                        └──────────────┘
```

### Allowed Transitions
- **Forward Progress**: Each phase must be completed before next
- **Revisions**: Can go back to previous phases for corrections
- **Approval Gates**: Plan→Task and Review→Complete require approval
- **Parallel Work**: Some phases allow concurrent activities

### State Persistence
All phase states and transitions are tracked in `.sdd/state.yaml`:

```yaml
project_id: sdd_1703123456
project_name: "My Awesome Project"
current_phase: execute
phases:
  specify:
    status: approved
    completed_at: "2024-01-01T10:00:00Z"
    agent_used: pm
  plan:
    status: approved
    completed_at: "2024-01-01T11:00:00Z"
    agent_used: architect
  task:
    status: approved
    completed_at: "2024-01-01T12:00:00Z"
    agent_used: developer
  execute:
    status: in_progress
    started_at: "2024-01-01T13:00:00Z"
    agent_used: developer
```

---

## Command Reference

### Core Workflow Commands
```bash
sdd init <name>           # Initialize project
sdd specify <desc>        # Generate specifications
sdd plan                  # Create architecture plan
sdd approve               # Approve current phase
sdd task                  # Break down into tasks
sdd execute               # Generate implementation guide
sdd review                # Quality assurance review
sdd status                # Show project status
```

### Workflow Management
```bash
sdd status                # Current phase and progress
sdd approve [comment]     # Approve phase with optional comment
sdd status --history      # Show phase transition history
```

### Quality Gates
- **Plan Approval**: Required before task breakdown
- **Review Approval**: Required before completion
- **Automatic Validation**: Basic checks on file existence and format

---

## Best Practices

### 1. Complete Each Phase Thoroughly
Don't rush through phases. Each phase builds quality into the next.

### 2. Use Approval Gates Effectively
Take approval seriously - it's your quality checkpoint.

### 3. Maintain Context Continuity
Each phase should reference and build upon previous phases.

### 4. Document Decisions
Use approval comments to document important decisions and rationale.

### 5. Iterate When Needed
It's okay to go back and revise earlier phases when new information emerges.

### 6. Parallel Development
Use task breakdowns to enable parallel development within teams.

### 7. Regular Status Checks
Use `sdd status` frequently to track progress and identify bottlenecks.

---

## Troubleshooting

### Stuck in Phase
**Issue:** Can't proceed to next phase

**Solutions:**
1. Check current phase status: `sdd status`
2. Ensure all required files exist
3. Check for approval requirements
4. Verify agent configurations

### Approval Rejected
**Issue:** Plan or review not approved

**Solutions:**
1. Review feedback and comments
2. Address identified issues
3. Revise and re-submit
4. Consider team consultation

### Context Issues
**Issue:** Agents not using proper context

**Solutions:**
1. Ensure LSP analysis completed
2. Check codebase structure
3. Verify agent configurations
4. Update context manually if needed

### State Corruption
**Issue:** Project state inconsistent

**Solutions:**
1. Check `.sdd/state.yaml` integrity
2. Reinitialize if necessary: `rm -rf .sdd && sdd init`
3. Restore from backup if available

---

The SDD workflow ensures that every feature goes through rigorous specification, design, implementation, and validation phases. This structured approach prevents common development pitfalls and ensures high-quality software delivery.