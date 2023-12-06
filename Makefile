SHELL := /bin/bash

# If the first argument is "run"...
ifeq (challenge,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif


.PHONY: challenge
challenge:
	./utils/make_challenge.sh ${RUN_ARGS}