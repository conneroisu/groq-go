library IEEE;

use IEEE.STD_LOGIC_1164.all;

entity FullAdder is

    port (
        i_A  : in  std_logic;
        i_B  : in  std_logic;
        Cin  : in  std_logic;
        Sum  : out std_logic;
        Cout : out std_logic
        );

end FullAdder;

architecture Structural of FullAdder is
    signal s_Adder, s_C1, s_C2 : std_logic;
begin

    s_Adder <= i_A xor i_B;
    s_C1    <= i_A and i_B;

    Sum  <= s_Adder xor Cin;
    s_C2 <= s_Adder and Cin;

    Cout <= s_C1 or s_C2;

end Structural;
