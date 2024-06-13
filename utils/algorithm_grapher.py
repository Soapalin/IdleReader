import matplotlib.pyplot as plt
from enum import Enum
import math

class IQTitle():
    EXTREMELY_LOW = 70
    VERY_LOW = 80
    LOW_AVERAGE = 90
    AVERAGE = 110
    HIGH_AVERAGE = 120
    VERY_HIGH = 130
    EXTREMELY_HIGH = 150
    SMARTY_PANTS = 200
    BIG_BRAIN = 400
    PLANETARY_BRAIN = 400
    GALAXY_BRAIN = 600


class FuncGraph:
    def __init__(self):
        pass

    def reading_increase(self):
        knowledge_cost = []
        knowledge_inc = []
        iq_required = []
        iq_gain = []
        reading_speed = []
        reading_speed_item = []
        knowledge_PerMin = []

        IQ = 40
        knowledgeCost = 100
        knowledgeIncrease = 5
        IQIncrease = 1
        IQRequired = 40
        readingSpeed = 1
        readingSpeed_item = 1
        pages = 100
        
        for i in range(600):
            IQ = IQ + 1
            if IQ <= 225:
                knowledgeCost = knowledgeCost*1.05 + 100

            auctionCost = (IQ/10)*knowledgeCost 
            knowledgeIncrease =  (IQ * pages)/10
            # knowledgeIncrease = knowledgeIncrease + (IQ * 10 + 20)/10
            if IQ <= 60:
                readingSpeed = 2
                readingSpeed_item = 2
            else:
                readingSpeed  =  2 + (readingSpeed+IQ)/IQ
                readingSpeed_item  = (readingSpeed + (readingSpeed+IQ)/IQ)*1.2

            print(f"================================================")
            print(f"knowledgeCost: {knowledgeCost} at book IQ: {IQ}")
            print(f"auctionCost: {auctionCost} at book IQ: {IQ}, reader IQ: {IQ}")

            print(f"knowledgeIncrease: {knowledgeIncrease} at book IQ: {IQ}")
            print(f"readingSpeed: {readingSpeed} at IQ: {IQ}")
            knowledgePerMin = ((readingSpeed * 60)/pages) * knowledgeIncrease
            print(f"knowledgePerMin: {knowledgePerMin}\n")
            print(f"Book read per min: {(readingSpeed * 60)/pages}\n")
            print(f"Time to Buy Book of this IQ (Seconds): {(knowledgeCost/(knowledgePerMin))*60}")
            print(f"Time to Buy Book of this IQ (Min): {knowledgeCost/(knowledgePerMin)}")
            print(f"Time to Buy Book of this IQ (Hours): {knowledgeCost//(knowledgePerMin*60)}")
            print(f"Time to Buy Book of this IQ (Days): {knowledgeCost//(knowledgePerMin*60*24)}\n\n")

            reading_speed.append(readingSpeed)
            knowledge_cost.append(knowledgeCost)
            knowledge_inc.append(knowledgeIncrease)
            reading_speed_item.append(readingSpeed_item)
            knowledge_PerMin.append(knowledgePerMin)

        x = [i for i in range(0, 600)]
        plt.plot(x, knowledge_cost,  label="knowledge_cost")
        plt.plot(x, knowledge_inc, label="knowledge_inc")
        plt.plot(x, reading_speed, label="reading_speed")
        plt.plot(x, reading_speed_item, label="reading_speed_item")
        plt.plot(x, knowledge_PerMin, label="knowledge_PerMin")
        plt.yscale("log")
        plt.legend()
        plt.title('Log values graph')
        plt.show()



    def main(self):
        self.reading_increase()


funcgraph = FuncGraph()
funcgraph.main()