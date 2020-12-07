# Helper script to generate certs from a csv
# Requires - run 'mockca generate -not-before x -not-after x' first to generate the ca
# Input - CSV file named 'identities.csv' with columns first,mi,last,edipi
# Output - to 'certs' directory crt, csr, key, and p12 files

import csv
import os

with open('identities.csv', 'r') as csvfile:
    reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in reader:
                print(row[0])
                os.system('openssl genrsa -out certs/'  + row[0] + '-' + row[1] + '-' + row[2] + '.key')
                os.system('openssl req -new -nodes -key certs/' + row[0] + '-' + row[1] + '-' + row[2] +  '.key -out certs/'  + row[0] + '-' + row[1] + '-' + row[2] + '.csr -subj \'/CN=' + row[2] + '.' + row[0] + '.' + row[1] + '.' + row[3] + '/\'')
                os.system('./mockca-build sign     -first-name "' + row[0] + '"     -last-name ' + row[2] + '     -middle-name ' + row[1] +'     -dod-id ' + row[3] + '     -email x+' + row[2] + '@dds.mil     -org marines     certs/'  + row[0] + '-' + row[1] + '-' + row[2] +  '.csr > certs/'  + row[0] + '-' + row[1] + '-' + row[2] +  '.crt')
                os.system('openssl pkcs12 -export -password pass:1234 -out certs/' + row[0] + '-' + row[1] + '-' + row[2] + '.p12 -inkey certs/'  + row[0] + '-' + row[1] + '-' + row[2] +  '.key  -in certs/'  + row[0] + '-' + row[1] + '-' + row[2] +  '.crt')
