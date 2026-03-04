import operator
import math
Lstr = input("Enter your function, like 1+x or 1 + x>>")
'''
#Lists
CLst = [x for x in Lstr]
nfull = []
ops_count = [0]*len(CLst)

#Dictionaries
ops = {
    '+': operator.add,
    '-': operator.sub,
    '*': operator.mul,
    '/': operator.truediv,
    '//': operator.floordiv,
    '%': operator.mod,
    '^': operator.pow,
    'sqrt': math.sqrt,
    'sin': math.sin,
    'cos': math.cos,
    'exp': math.exp,
    'log': math.log
}

i = 0
nchar = ''
while i < len(CLst):
    if CLst[i].isspace():
        i += 1
        continue
    if CLst[i].isdigit():
        nchar += CLst[i]
        try:
            if CLst[i+1].isdigit() == False:
                nfull.append(nchar)
                nchar = ''
        except:
            nfull.append(nchar)
    elif ''.join(CLst[i:i+4]) in ops:
        nfull.append(''.join(CLst[i:i+4]))
        ops_count[i] = 1
        i += 3
    elif ''.join(CLst[i:i+3]) in ops:
        nfull.append(''.join(CLst[i:i+3]))
        ops_count[i] = 1
        i += 2
    elif ''.join(CLst[i:i+2]) in ops:
        nfull.append(''.join(CLst[i:i+2]))
        ops_count[i] = 1
        i += 1
    elif CLst[i] in ops:
        nfull.append(CLst[i])
        ops_count[i] = 1
    elif CLst[i].isalpha():
        nfull.append(CLst[i])
    i+=1
nfullm = []
for n in nfull:
    if n.isdigit():
        nfullm.append(float(n))
    else:
        nfullm.append(n)
print(nfullm, ops_count)
"""
result = st = 0
for op in range(len(nfull)):
    if st == ops_count:
        break
    if nfull[op] in ops:
        try:
            if len(nfull[op]) <= 2:
                result += ops[nfull[op]](nfullm[op-1], nfullm[op+1])
            else:
                result += ops[nfull[op]](nfullm[op+1])
            st += 1
        except:
            print("!!!_ERROR_!!!")
    
print(f"Final list>>{result}")
"""
def steps(list_with_separated_elements, st, final_st, op_or_not, result):
    print(list_with_separated_elements)
    if st >= final_st:
        return result
    else:
        for i in range(len(list_with_separated_elements)):
            if op_or_not[i] == 1:
                if len(list_with_separated_elements[i]) <= 2:
                    print("Step of 2 el summ<>", st+1)
                    result = ops[list_with_separated_elements[i]](result, list_with_separated_elements[i+1])
                    return steps(list_with_separated_elements, st+1, final_st, op_or_not, 
                            result)
                else:
                    return steps(list_with_separated_elements, st+1, final_st, op_or_not,
                            ops[list_with_separated_elements[i]](list_with_separated_elements[i+1]))
                
                


#Отдельные эл, стартовое ч, конечное ч, список с проверкой на оператор, результат
print(f"Result >> {steps(nfullm, 0, sum(ops_count), ops_count, 0)}")'''

def diсhotomy():
    
    return 0

print(eval(Lstr))
