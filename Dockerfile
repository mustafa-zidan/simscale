FROM alpine

ARG INPUT_FILE="input.txt"
ARG OUTPUT_FILE="out.txt"

ADD bin/simscale_*_linux_amd64 /simscale

CMD /simscale --in-file=${INPUT_FILE} --out-file=${OUTPUT_FILE}