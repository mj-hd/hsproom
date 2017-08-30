#!/usr/bin/python

import time
import sys
import urllib.request
import urllib.parse
import os
import json
import re
from pyknp import Jumanpp
from gensim import models
from gensim.models.doc2vec import LabeledSentence

TickInterval    = 1
LearnInterval   = 60 * 60 * 60 * 24
ConvertInterval = 3
FailureInterval = 10
HsproomProtocol = "http"
HsproomHost     = "127.0.0.1:8080"
HsproomApiRoot  = "/api/"
CorpusSize      = 100
CorpusEpoch     = 30
VectorSize      = 512
ModelPath       = "./learner/doc2vec.model"

def compose_url(endpoint):
    return HsproomProtocol + "://" + HsproomHost + HsproomApiRoot + endpoint

def post(endpoint, params):
    url = compose_url(endpoint)
    data = urllib.parse.urlencode(params).encode("utf8")
    res = urllib.request.urlopen(url, data=data)
    r = res.read().decode("utf8")
    j = json.loads(r)
    return j

def get(endpoint, params):
    url = compose_url(endpoint)
    url += "?{0}".format(urllib.parse.urlencode(params))
    res = urllib.request.urlopen(url)
    r = res.read().decode("utf8")
    j = json.loads(r)
    return j

    time.sleep(TickInterval)
def get_corpus():
    params = {"n": CorpusSize}
    res = get("batch/get/corpus/", params)
    return res["Corpus"]

def get_unprocessed_document():
    res = get("batch/get/document/unprocessed/", {})
    return res["ID"], res["Document"]

def post_vector(id, vector):
    params = {"p": id, "v": ",".join([str(v) for v in vector])}
    res = post("batch/save/vector/", params)

def split_into_words(text):
    result = Jumanpp().analysis(text)
    return [mrph.midasi for mrph in result.mrph_list()]

def normalize_doc(text):
    text = re.sub(r'[!-~]', "", text)#半角記号,数字,英字
    text = re.sub(r'[︰-＠]', "", text)#全角記号
    text = re.sub('[\n\r]', " ", text)#改行文字
    text = re.sub('\s+', " ", text)
    print("\n[" + text + "]\n")
    return text

def doc_to_sentence(doc, name):
    doc = normalize_doc(doc)
    if not doc:
        words = []
    else:
        words = split_into_words(doc)
    return LabeledSentence(words=words, tags=[name])

def corpus_to_sentences(corpus):
    for i, doc in enumerate(corpus):
        if not doc:
            continue
        time.sleep(1)
        sys.stdout.write("* 前処理中 {}/{}".format(i, len(corpus)))
        yield doc_to_sentence(doc, "CORPUS-" + str(i))

def learn():
    corpus = get_corpus()
    sentences = corpus_to_sentences(corpus)
    model = models.Doc2Vec(sentences, dm=0, size=VectorSize, window=15, alpha=.025, min_alpha=.025, min_count=1, sample=1e-6)
    for epoch in range(CorpusEpoch):
        print("Epoch: {}".format(epoch + 1))
        model.train(sentences, total_examples=model.corpus_count, epochs=model.iter)
        model.alpha -= (0.025 - 0.0001) / 19
        model.min_alpha = model.alpha

    model.save(ModelPath)
    print("訓練終了:" + ModelPath)
    return model

def load_model():
    if not os.path.exists(ModelPath):
        print("モデルが存在しない->訓練")
        return learn()

    return models.Doc2Vec.load(ModelPath)

def convert_to_vector(model):
    id, doc = get_unprocessed_document()

    doc = normalize_doc(doc)
    if not doc:
        convert_to_vector(model)
        return

    print("変換開始...")
    words = split_into_words(doc)
    vector = model.infer_vector(words)
    post_vector(id, vector)
    print("変換完了:" + str(id))

print("モデル読み込み...")
model = load_model()

print("待機...")
healthy = True
timer   = 1
while healthy:
    try:
        if timer % LearnInterval == 0:
            print("学習開始...")
            model = learn()

        if timer % ConvertInterval == 0:
            convert_to_vector(model)
    except:
        time.sleep(FailureInterval)

    time.sleep(TickInterval)
    timer += 1
